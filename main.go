package main

import (
  "fmt"
  "flag"
  "net/http"
  "time"
)

const MAX_TWEETS = 1000
const REFRESH_INTERVAL = 60 * time.Second

var (
  consumerKey    string
  consumerSecret string
  
  keyword        string
  reservoir      = NewReservoirSampler(MAX_TWEETS, NewPseudoRangomNumberGenerator())
  globalStats    []wordStat
)

func initApp() {
  flag.StringVar(&consumerKey, "ck", "", "Consumer Key")
  flag.StringVar(&consumerSecret, "cs", "", "Consumer Secret")
  flag.Parse()
}

func showStats() {
  c := time.Tick(REFRESH_INTERVAL)
  for {
    <-c
    samples := reservoir.GetSamples()
    reservoir = NewReservoirSampler(MAX_TWEETS, NewPseudoRangomNumberGenerator())
    tweets := make([]*Tweet, 0, len(samples))
    for _, sample := range samples {
      tweets = append(tweets, sample.(*Tweet))
    }
    globalStats = getStats(tweets)
    fmt.Println(len(samples), globalStats)
  }
}

func showAllStats(w http.ResponseWriter, r *http.Request) {
  if(globalStats == nil) {
    fmt.Fprintln(w, "no stats")
    return
  }
  fmt.Fprintf(w, "<table>")
  for _, stat := range globalStats {
    wordDetails := WordIndex[stat.word]
    fmt.Fprintf(w, "<tr><td style='background-color:%v'>&nbsp;</td><td>%v</td><td>%v</td></tr>", wordDetails.Color(), stat.word, stat.count)
  }
  fmt.Fprintf(w, "</table>")  
}

func showTopWord(w http.ResponseWriter, r *http.Request) {
  if(globalStats == nil || len(globalStats) < 1) {
    w.WriteHeader(412)
    fmt.Fprintln(w, "no stats")
    return
  }

  winner := globalStats[0]
  wd := WordIndex[winner.word]
  fmt.Fprintf(w, "<Statuses><Status>")
  fmt.Fprintf(w, "<Title>%v</Title>", winner.word)
  fmt.Fprintf(w, "<Colour>%v,%v,%v</Colour>", wd.Red(), wd.Green(), wd.Blue())
  fmt.Fprintf(w, "<Red>%v</Red><Green>%v</Green><Blue>%v</Blue>", wd.Red(), wd.Green(), wd.Blue())
  fmt.Fprintf(w, "<Hex>%v</Hex>", wd.Color())
  fmt.Fprintf(w, "<RGBA>%v,%v,%v,%v</RGBA>", wd.Red(), wd.Green(), wd.Blue(), wd.Alpha())
  fmt.Fprintf(w, "</Status></Statuses>")
  
}

func serveStats() {
  http.HandleFunc("/", showTopWord)
  http.HandleFunc("/stats", showAllStats)
  http.ListenAndServe(":9080", nil)
}


func main() {
  initApp()

  go showStats()
  go serveStats()

  var tweets chan *Tweet
  tweets, err := GetTweets(consumerKey, consumerSecret)
  if err != nil {
    fmt.Println("Error:", err)
    return
  }

  for tweet := range tweets {
    if tweet == nil {
      break
    }
    if tweet.Id == 0 {
      continue
    }
    reservoir.Add(tweet)
  }

}
