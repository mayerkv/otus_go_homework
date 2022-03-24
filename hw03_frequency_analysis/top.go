package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type wordCount struct {
	w   string
	cnt int
}

func Top10(s string) []string {
	words := strings.Fields(s)
	if len(words) == 0 {
		return nil
	}

	wordsCountMap := map[string]int{}
	for _, word := range words {
		wordsCountMap[word]++
	}

	wordsCountSlice := make([]wordCount, 0, len(wordsCountMap))

	for w, cnt := range wordsCountMap {
		wordsCountSlice = append(wordsCountSlice, wordCount{w, cnt})
	}

	sort.Slice(wordsCountSlice, func(i, j int) bool {
		a, b := wordsCountSlice[i], wordsCountSlice[j]

		if a.cnt == b.cnt {
			return strings.Compare(a.w, b.w) < 0
		}

		return a.cnt > b.cnt
	})

	wordsCount := 10
	if len(wordsCountSlice) < wordsCount {
		wordsCount = len(wordsCountSlice)
	}

	res := make([]string, 0, wordsCount)
	for _, item := range wordsCountSlice[:wordsCount] {
		res = append(res, item.w)
	}

	return res
}
