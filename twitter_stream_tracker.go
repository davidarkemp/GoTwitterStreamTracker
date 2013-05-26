package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mrjones/oauth"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const TOKEN_FILE = ".twitter_oauth"

var consumer       *oauth.Consumer

type User struct {
	Id          int64
	Screen_Name string
}

type Tweet struct {
	Id   int64
	Text string
	User User
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

func getAccessToken() (accessToken *oauth.AccessToken, err error) {
	accessToken = getSavedAccessToken()

	if accessToken != nil {
		return
	}

	accessToken, err = authApp()

	if err != nil {
		return
	}

	if accessToken == nil {
		err = errors.New("cannot get Error token")
		return
	}

	file, err := os.Create(TOKEN_FILE)

	if err != nil {
		return
	}

	defer file.Close()
	file.WriteString(accessToken.Token + "%" + accessToken.Secret)

	return
}

func GetTweets(consumerKey, consumerSecret string) (tweets chan *Tweet, err error) {
	consumer = oauth.NewConsumer(
    consumerKey,
    consumerSecret,
    oauth.ServiceProvider{
      RequestTokenUrl:   "http://api.twitter.com/oauth/request_token",
      AuthorizeTokenUrl: "https://api.twitter.com/oauth/authorize",
      AccessTokenUrl:    "https://api.twitter.com/oauth/access_token",
    })

	var accessToken *oauth.AccessToken

	accessToken, err = getAccessToken()
	if err != nil {
		return nil, err
	}

	/*
		keywords := make([]string, len(WordList))
		for _, word := range WordList {
			keywords = append(keywords, word.Word)
		}
		keyword = strings.Join(keywords, ",")
		result, err := consumer.Post("https://stream.twitter.com/1.1/statuses/filter.json", map[string]string{"track": keyword}, accessToken)
	*/
	result, err := consumer.Get("https://stream.twitter.com/1.1/statuses/sample.json", nil, accessToken)

	if err != nil {
		return
	}

	tweets = make(chan *Tweet)
	go readStream(result.Body, tweets)
	return
}

