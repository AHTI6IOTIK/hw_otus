package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type WordRepeat struct {
	value string
	count int
}

type WordRepetitionCounter struct {
	convertedCache []WordRepeat
	words          map[string]int
	step           int
}

func New(size int) WordRepetitionCounter {
	return WordRepetitionCounter{
		step:  0,
		words: make(map[string]int, size),
	}
}

func (a *WordRepetitionCounter) GetWords() []WordRepeat {
	if len(a.convertedCache) > 0 {
		return a.convertedCache
	}

	result := make([]WordRepeat, 0, len(a.words))
	for k, v := range a.words {
		result = append(result, WordRepeat{
			value: k,
			count: v,
		})
	}

	a.convertedCache = result

	return result
}

func (a *WordRepetitionCounter) GetStep() int {
	return a.step
}

func (a *WordRepetitionCounter) Append(word string) {
	if v, ok := a.words[word]; ok {
		a.words[word] = v + 1
	} else {
		a.words[word] = 1
	}
}

func (a *WordRepetitionCounter) GetTop(count int) []WordRepeat {
	wordsCount := len(a.words)
	if count > wordsCount {
		count = wordsCount
	}

	takenItems := make(map[string]struct{}, a.step)
	topItems := make([]string, 0, count)

	for len(topItems) < count {
		maxItem := 0
		selectedItem := ""

		for i, v := range a.words {
			_, isAlreadySelected := takenItems[i]
			if v > maxItem && !isAlreadySelected {
				maxItem = v
				selectedItem = i
			}
		}

		topItems = append(topItems, selectedItem)
		takenItems[selectedItem] = struct{}{}
	}

	result := make([]WordRepeat, 0, len(topItems))
	for _, v := range topItems {
		result = append(result, WordRepeat{
			value: v,
			count: a.words[v],
		})
	}

	return result
}

func Top10(str string) []string {
	if str == "" {
		return []string{}
	}

	splitStr := strings.Fields(str)
	wordCounter := New(len(splitStr))

	for _, item := range splitStr {
		wordCounter.Append(item)
	}

	top10 := wordCounter.GetTop(10)

	sort.Slice(top10, func(i, j int) bool {
		if top10[i].count == top10[j].count {
			return top10[i].value < top10[j].value
		}

		return top10[i].count > top10[j].count
	})

	result := make([]string, 0, len(top10))

	for _, v := range top10 {
		result = append(result, v.value)
	}

	return result
}
