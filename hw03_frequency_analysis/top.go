package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Top10(s string) []string {
	type counter struct {
		count int64
		pos   int
	}

	words := strings.Fields(s)
	freq := make(map[string]counter)
	for i, v := range words {
		freq[v] = counter{
			count: freq[v].count + 1,
			pos:   i,
		}
	}

	freqSlice := make([]counter, 0, len(freq))
	for _, v := range freq {
		freqSlice = append(freqSlice, v)
	}
	sort.Slice(freqSlice, func(i, j int) bool {
		return freqSlice[j].count < freqSlice[i].count ||
			freqSlice[j].count == freqSlice[i].count && words[freqSlice[j].pos] > words[freqSlice[i].pos]
	})

	length := min(10, len(freqSlice))
	res := make([]string, 0, length)
	for i := 0; i < length; i++ {
		res = append(res, words[freqSlice[i].pos])
	}
	return res
}
