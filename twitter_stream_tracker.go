package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/mrjones/oauth"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
	"time"
  "regexp"
)

const TOKEN_FILE = ".twitter_oauth"
const MAX_TWEETS = 10000

var (
	consumerKey    string
	consumerSecret string
	consumer       *oauth.Consumer
	keyword        string
	stopWords      = map[string]bool{
		"the":  true,
		"rt":   true,
		"a":    true,
		"to":   true,
		"me":   true,
		"my":   true,
		"and":  true,
		"in":   true,
		"is":   true,
		"of":   true,
		"so":   true,
		"for":  true,
		"on":   true,
		"at":   true,
		"this": true,
		"it":   true,
		"with": true,
		"that": true,
  }
  incWords = map[string]*regexp.Regexp {
    "love": regexp.MustCompile(`\blove\b`),
    "hate": regexp.MustCompile(`\bhate\b`),
    "car": regexp.MustCompile(`\bcar\b`),
    "home": regexp.MustCompile(`\bhome\b`),
  }
)

type User struct {
	Id          int64
	Screen_Name string
}

type Tweet struct {
	Id    int64
	Text  string
	User  User
	Words map[string]int64
}

func initApp() {
	flag.StringVar(&consumerKey, "ck", "", "Consumer Key")
	flag.StringVar(&consumerSecret, "cs", "", "Consumer Secret")
	flag.StringVar(&keyword, "keyword", "London", "search word")
	flag.Parse()

	consumer = oauth.NewConsumer(
		consumerKey,
		consumerSecret,
		oauth.ServiceProvider{
			RequestTokenUrl:   "http://api.twitter.com/oauth/request_token",
			AuthorizeTokenUrl: "https://api.twitter.com/oauth/authorize",
			AccessTokenUrl:    "https://api.twitter.com/oauth/access_token",
		})
}

func getSavedAccessToken() *oauth.AccessToken {
	var token oauth.AccessToken
	file, err := os.Open(TOKEN_FILE)

	if err == nil {
		defer file.Close()
		bytes, err := ioutil.ReadAll(file)

		if err == nil {
			line := string(bytes)

			token_strs := strings.Split(line, "%")
			token.Token = token_strs[0]
			token.Secret = token_strs[1]

			return &token
		}
	}

	return nil
}

func authApp() (*oauth.AccessToken, error) {
	requestToken, url, err := consumer.GetRequestTokenAndUrl("oob")
	if err != nil {
		log.Fatal(err)

		return nil, err
	}

	fmt.Println("(1) Go to: " + url)
	fmt.Println("(2) Grant access, you should get back a verification code.")
	fmt.Println("(3) Enter that verification code here: ")

	verificationCode := ""
	fmt.Scanln(&verificationCode)

	accessToken, err := consumer.AuthorizeToken(requestToken, verificationCode)
	if err != nil {
		return nil, err
	}

	return accessToken, err
}

func readStream(reader io.Reader, ch chan *Tweet) {
	dec := json.NewDecoder(reader)
	for {
		var t Tweet
		err := dec.Decode(&t)

		if err != nil {
			ch <- nil
			fmt.Println(err)
			break
		}

		ch <- &t
	}
}

func dealWithTweet(p *Tweet, c chan *Tweet) {
	parts := strings.Split(strings.ToLower(p.Text), " ")
	p.Words = make(map[string]int64)
	for _, word := range parts {
		if len(word) < 2 || word[0] == '@' {
			continue
		}
		if count, ok := p.Words[word]; ok {
			p.Words[word] = count + 1
		} else {
			p.Words[word] = 1
		}
	}

	c <- p

}

func showTweets(output chan *Tweet) {
	for {
		<-output
		//fmt.Printf("%v:%v:%v\n", p.Id, p.Text, p.Words)
	}
}

type wordStat struct {
	word  string
	count int64
}

type wordStats []wordStat

func (w wordStats) Len() int {
	return len(w)
}

func (w wordStats) Swap(i, j int) {
	w[i], w[j] = w[j], w[i]
}

type ByCount struct {
	wordStats
}

func (w ByCount) Less(i, j int) bool {
	return w.wordStats[i].count < w.wordStats[j].count
}

func showStats(r *Reservoir) {
	c := time.Tick(5 * time.Second)
	for {
		<-c
		samples := r.GetSamples()
		words := make(map[string]wordStat, len(incWords))
		for _, sample := range samples {
			p := sample.(*Tweet)

      for word, regex := range incWords {
        count := int64(len(regex.FindAllString(p.Text, -1)))
      	existingCount, _ := words[word]
				words[word] = wordStat{word, existingCount.count + count}
			}
		}
		stats := make([]wordStat, len(words))
		for _, stat := range words {
			stats = append(stats, stat)
		}
		sort.Sort(sort.Reverse(ByCount{stats}))

		fmt.Println(len(samples), stats)

	}
}

func main() {
	initApp()
	var accessToken *oauth.AccessToken

	accessToken = getSavedAccessToken()

	reservoir := NewReservoirSampler(MAX_TWEETS, NewPseudoRangomNumberGenerator())

	if accessToken == nil {
		var err error
		accessToken, err = authApp()
		//Save token

		if err == nil {
			file, err := os.Create(TOKEN_FILE)

			if err == nil {
				defer file.Close()
				file.WriteString(accessToken.Token + "%" + accessToken.Secret)
			}
		} else {
			fmt.Println(err)
			os.Exit(0)
		}
	}

	if accessToken != nil {

		//result, err := consumer.Post("https://stream.twitter.com/1.1/statuses/filter.json", map[string]string{"track": keyword}, accessToken)
    result, err := consumer.Get("https://stream.twitter.com/1.1/statuses/sample.json", nil, accessToken)

		if err != nil {
      fmt.Println("Error: %v", err)
      return
    }
		
    ch := make(chan *Tweet)
		output := make(chan *Tweet)
		go showStats(reservoir)
		go showTweets(output)
		go readStream(result.Body, ch)

		for {
			p := <-ch
			if p == nil {
				break
			}
			if p.Id == 0 {
				continue
			}

			if reservoir.Add(p) {
				go dealWithTweet(p, output)
			}
		}
		
	}
}
