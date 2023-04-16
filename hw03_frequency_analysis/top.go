package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

type wordCounter struct {
	Word  string
	Count int64
}

var rx = *regexp.MustCompile(`([\p{L}\p{N}_-]*)([^-?.,:;!'"#@/|+=#$%^&*(){}\[\]])`)

func Top10(text string) []string {
	words := strings.Fields(strings.ToLower(text))

	if len(words) == 0 {
		return nil
	}

	sort.Strings(words)

	freqWords := []wordCounter{}

	for _, elem := range words {
		cleanWord := rx.FindString(elem)
		if cleanWord == "" {
			continue
		}
		elemIndex := sort.Search(len(freqWords), func(i int) bool {
			return freqWords[i].Word == cleanWord
		})
		if elemIndex < len(freqWords) {
			freqWords[elemIndex].Count++
		} else {
			freqWords = append(freqWords, wordCounter{cleanWord, 1})
		}
	}

	sort.SliceStable(freqWords, func(i, j int) bool {
		return freqWords[i].Count > freqWords[j].Count
	})

	var topTenWords []wordCounter
	if len(freqWords) < 10 {
		ln := len(freqWords)
		topTenWords = freqWords[:ln]
	} else {
		topTenWords = freqWords[:10]
	}

	resultWords := make([]string, 0)
	for _, elem := range topTenWords {
		resultWords = append(resultWords, elem.Word)
	}

	return resultWords
}
