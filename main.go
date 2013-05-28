package main

import (
  "fmt"
  "flag"
  "net/http"
  "log"
  "time"
  "container/list"
)

const MAX_TWEETS = 1000
const REFRESH_INTERVAL = 60 * time.Second
const MAX_HISTORY = 100

var (
  consumerKey    string
  consumerSecret string
  globalStats    wordStats
  keyword        string
  reservoir      = NewReservoirSampler(MAX_TWEETS, NewPseudoRangomNumberGenerator())
  wordHistory    = make(map[string]*list.List)
  numberOfSamples int64
  wordTotals     = make(map[string]int64)
)

func initApp() {
  flag.StringVar(&consumerKey, "ck", "", "Consumer Key")
  flag.StringVar(&consumerSecret, "cs", "", "Consumer Secret")
  flag.Parse()
}

func drainReservoir() (tweets []*Tweet) {
  samples := reservoir.GetSamples()
  if(len(samples) == 0) {
      return
  }
  reservoir = NewReservoirSampler(MAX_TWEETS, NewPseudoRangomNumberGenerator())
  tweets = make([]*Tweet, 0, len(samples))
  for _, sample := range samples {
    tweets = append(tweets, sample.(*Tweet))
  }
  return
}

func showStats() {
  c := time.Tick(REFRESH_INTERVAL)
  for {
    <-c
    tweets := drainReservoir() 
    if(tweets == nil) {
      continue
    } 
    numberOfSamples += 1
    log.Println("updated stats", len(tweets))
    stats, raw := getStats(tweets)
    globalStats = stats
    
    for word, count := range raw {
      if wordHistory[word] == nil {
        wordHistory[word] = list.New()
      }
      wordHistory[word].PushFront(count)
      for wordHistory[word].Len() > MAX_HISTORY {
        wordHistory[word].Remove(wordHistory[word].Back())
      }
      wordTotals[word] = wordTotals[word] + int64(count)
    }
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
    fmt.Fprintf(w, "<tr><td style='background-color:%v'>&nbsp;</td><td>%v</td>", wordDetails.Color(), stat.word)

    var average float64
    if(numberOfSamples > 0) {
      average = float64(wordTotals[stat.word]) / float64(numberOfSamples)
    }
    fmt.Fprintf(w, "<td>w:%f<td><td>avg:%f<td>", wordDetails.Weight(), average)

    for history:= wordHistory[stat.word].Front(); history != nil; history = history.Next() {
      fmt.Fprintf(w, "<td>%d<td>",history.Value)
    }
    fmt.Fprintf(w, "</tr>")
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
  if err := http.ListenAndServe(":9080", nil); err != nil {
    log.Fatal("Web server", err)
  }
}


func main() {
  initApp()

  go showStats()
  go serveStats()

  timeout_duration := 10 * time.Second
  timer := time.NewTimer(timeout_duration);
  for{
    var tweets chan *Tweet
    tweets, err := GetTweets(consumerKey, consumerSecret)
    if err != nil {
      log.Println("Error:", err)
    } else {  
      for {
        select {
          case tweet := <- tweets: 
            if tweet == nil {
              break
            }
            if tweet.Id == 0 {
              continue
            }
            reservoir.Add(tweet)
          
          case <- timer.C:
            log.Println("disconnected")
            break
        }
        timer.Reset(timeout_duration)
      }
    }
    time.Sleep(5 * time.Second)
  }


}
