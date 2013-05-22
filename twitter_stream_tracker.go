package main

import (
	"encoding/json"
	"fmt"
	"github.com/mrjones/oauth"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
    "flag"
)

const TOKEN_FILE = ".twitter_oauth"

var consumerKey string
var consumerSecret string
var consumer *oauth.Consumer
var keyword string

type User struct {
    Id  int64
    Screen_Name string
}

type Tweet struct {
	Id   int64
	Text string
    User User
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

func main() {
	initApp()
	var accessToken *oauth.AccessToken

	accessToken = getSavedAccessToken()

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

		result, err := consumer.Post("https://stream.twitter.com/1.1/statuses/filter.json", map[string]string{"track": keyword}, accessToken)

		if err == nil {
			ch := make(chan *Tweet)
			go readStream(result.Body, ch)

			for {
				p := <-ch

				if p == nil {
					break
				}
                if(p.Id == 0) {
                    continue
                }
				fmt.Print(p.Id)
                fmt.Print(":", p.User.Screen_Name)
                parts := strings.Split(strings.ToLower(p.Text), " ")
                counts := make(map[string] int64)
                for _, word := range parts {
                    if len(word) < 2 || word[0] == '@' {
                        continue
                    }
                    if count, ok := counts[word]; ok {
                        counts[word] = count + 1
                    } else {
                        counts[word] = 1
                    }
                }
                fmt.Print(":",p.Text)
				fmt.Print(":", counts)
                fmt.Println("")
			}
		} else {
            fmt.Println("Error: %v", err)
        }
	}
}
