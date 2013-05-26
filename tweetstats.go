package main

import (
	"sort"
)

type wordStat struct {
	word  string
	count int
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

func wordCounts(p *Tweet) map[string]int {
	words := make(map[string]int)
	for _, wordDetails := range WordList {
		words[wordDetails.Word()] = len(wordDetails.Regexp().FindAllString(p.Text, -1))
	}
	return words
}

func getStats(tweets []*Tweet) (stats []wordStat) {
	words := make(map[string]int)
	for _, tweet := range tweets {

		tweetWords := wordCounts(tweet)

		for word, count := range tweetWords {
			if len(word) == 0 {
				continue
			}
			words[word] += count
		}
	}

	stats = make([]wordStat, 0, len(words))
	for word, count := range words {
		if len(word) == 0 {
			continue
		}
		weight := WordIndex[word].Weight()

		stats = append(stats, wordStat{word: word, count: int(float32(count) * weight)})
	}
	sort.Sort(sort.Reverse(ByCount{stats}))
	return
}
