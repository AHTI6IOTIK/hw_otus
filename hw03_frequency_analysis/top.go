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
	words []WordRepeat
	step  int
}

func New(size int) WordRepetitionCounter {
	return WordRepetitionCounter{
		step:  0,
		words: make([]WordRepeat, size),
	}
}

func (a *WordRepetitionCounter) GetWords() []WordRepeat {
	return a.words
}

func (a *WordRepetitionCounter) GetStep() int {
	return a.step
}

func (a *WordRepetitionCounter) Append(word string) {
	for i, v := range a.words {
		if v.value == word {
			v.count++
			a.words[i] = v

			return
		}
	}

	a.words[a.step] = WordRepeat{
		value: word,
		count: 1,
	}

	a.step++
}

func (a *WordRepetitionCounter) GetTop(count int) []WordRepeat {
	if count > a.step {
		count = a.step
	}

	takenItems := make(map[int]struct{}, a.step)
	topItems := make([]WordRepeat, 0, count)

	for len(topItems) < count {
		selectedValue := 0
		maxItem := WordRepeat{}
		for i := 0; i < a.step; i++ {
			_, alreadySelected := takenItems[i]
			currentItem := a.words[i]
			if currentItem.count > maxItem.count && !alreadySelected {
				maxItem = currentItem
				selectedValue = i
			}
		}

		topItems = append(topItems, a.words[selectedValue])
		takenItems[selectedValue] = struct{}{}
	}

	return topItems
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
