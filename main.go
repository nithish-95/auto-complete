package main

import (
	"fmt"
	"runtime"
	"sort"
	"time"
)

// -----------------------------------------
// Algorithm_1: Contextual Bigram-Based Trie
// -----------------------------------------

type TrieNodeA1 struct {
	children  map[rune]*TrieNodeA1
	isEnd     bool
	frequency int
}

type TrieA1 struct {
	root        *TrieNodeA1
	bigramTable map[string]map[string]int
}

func NewTrieNodeA1() *TrieNodeA1 {
	return &TrieNodeA1{children: make(map[rune]*TrieNodeA1)}
}

func NewTrieA1() *TrieA1 {
	return &TrieA1{
		root:        NewTrieNodeA1(),
		bigramTable: make(map[string]map[string]int),
	}
}

func (t *TrieA1) Insert(word string) {
	node := t.root
	for _, char := range word {
		if _, exists := node.children[char]; !exists {
			node.children[char] = NewTrieNodeA1()
		}
		node = node.children[char]
	}
	node.isEnd = true
	node.frequency++
}

func (t *TrieA1) BuildBigramTable(corpus []string) {
	for i := 0; i < len(corpus)-1; i++ {
		word1 := corpus[i]
		word2 := corpus[i+1]

		if _, exists := t.bigramTable[word1]; !exists {
			t.bigramTable[word1] = map[string]int{"_total": 0}
		}
		t.bigramTable[word1][word2]++
		t.bigramTable[word1]["_total"]++
	}
}

func (t *TrieA1) searchPrefix(prefix string) *TrieNodeA1 {
	node := t.root
	for _, char := range prefix {
		if child, exists := node.children[char]; exists {
			node = child
		} else {
			return nil
		}
	}
	return node
}

func (t *TrieA1) collectCompletions(node *TrieNodeA1, prefix string) []struct {
	word      string
	frequency int
} {
	var results []struct {
		word      string
		frequency int
	}

	var dfs func(*TrieNodeA1, []rune)
	dfs = func(currentNode *TrieNodeA1, path []rune) {
		if currentNode.isEnd {
			results = append(results, struct {
				word      string
				frequency int
			}{word: string(path), frequency: currentNode.frequency})
		}
		for char, childNode := range currentNode.children {
			dfs(childNode, append(path, char))
		}
	}

	dfs(node, []rune(prefix))
	return results
}

func (t *TrieA1) rankByContextualProbability(prefix string, completions []struct {
	word      string
	frequency int
}) []struct {
	word        string
	probability float64
} {
	if contextData, exists := t.bigramTable[prefix]; exists {
		totalFrequency := contextData["_total"]

		var ranked []struct {
			word        string
			probability float64
		}
		for _, completion := range completions {
			bigramFreq := contextData[completion.word]
			probability := float64(bigramFreq) / float64(totalFrequency)
			ranked = append(ranked, struct {
				word        string
				probability float64
			}{word: completion.word, probability: probability})
		}

		sort.Slice(ranked, func(i, j int) bool {
			return ranked[i].probability > ranked[j].probability
		})
		return ranked
	}

	// If no context is available, use frequency
	totalFreq := 0
	for _, completion := range completions {
		totalFreq += completion.frequency
	}

	var ranked []struct {
		word        string
		probability float64
	}
	for _, completion := range completions {
		probability := float64(completion.frequency) / float64(totalFreq)
		ranked = append(ranked, struct {
			word        string
			probability float64
		}{word: completion.word, probability: probability})
	}

	sort.Slice(ranked, func(i, j int) bool {
		return ranked[i].probability > ranked[j].probability
	})
	return ranked
}

func (t *TrieA1) Autocomplete(prefix string, k int) []struct {
	word        string
	probability float64
} {
	node := t.searchPrefix(prefix)
	if node == nil {
		return nil
	}

	completions := t.collectCompletions(node, prefix)
	rankedCompletions := t.rankByContextualProbability(prefix, completions)

	if k > len(rankedCompletions) {
		k = len(rankedCompletions)
	}
	return rankedCompletions[:k]
}

// -----------------------------------------
// Algorithm_2: Frequency-Based Trie
// -----------------------------------------

type NodeA2 struct {
	children    map[rune]*NodeA2
	isEndOfWord bool
	frequency   int
}

type TriesA2 struct {
	root *NodeA2
}

func initTriesA2() *TriesA2 {
	return &TriesA2{
		root: &NodeA2{
			isEndOfWord: false,
			children:    make(map[rune]*NodeA2),
			frequency:   0,
		},
	}
}

func (t *TriesA2) Insert(word string) {
	current := t.root
	for _, char := range word {
		node, ok := current.children[char]
		if !ok {
			node = &NodeA2{
				isEndOfWord: false,
				children:    make(map[rune]*NodeA2),
				frequency:   0,
			}
			current.children[char] = node
		}
		current = node
	}
	current.isEndOfWord = true
	current.frequency++
}

func (t *TriesA2) getFrequency(word string) int {
	current := t.root
	for _, char := range word {
		node, ok := current.children[char]
		if !ok {
			return 0
		}
		current = node
	}
	return current.frequency
}

func (t *TriesA2) Autocomplete(prefix string) []string {
	current := t.root
	for _, char := range prefix {
		node, ok := current.children[char]
		if !ok {
			return []string{}
		}
		current = node
	}

	var results []string
	collectWordsA2(current, prefix, &results)

	// Sort by frequency
	sort.Slice(results, func(i, j int) bool {
		return t.getFrequency(results[i]) > t.getFrequency(results[j])
	})

	// Limit to top 10 suggestions
	if len(results) > 10 {
		results = results[:10]
	}

	return results
}

func collectWordsA2(node *NodeA2, prefix string, results *[]string) {
	if node.isEndOfWord {
		*results = append(*results, prefix)
	}
	for char, child := range node.children {
		collectWordsA2(child, prefix+string(char), results)
	}
}

// -----------------------------------------
// Metrics and Evaluation Code
// -----------------------------------------

// Measures memory usage and returns bytes allocated
func getMemoryUsage() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.Alloc
}

// Measures suggestion quality by comparing returned suggestions to an ideal set
// Here, we define "ideal" as a set of words we expect to see at the top.
// We measure how many top expected words are present in the suggestions.
func measureSuggestionQuality(got []string, ideal []string) float64 {
	if len(ideal) == 0 {
		return 1.0 // If no ideal set given, can't measure quality – assume perfect
	}
	hitCount := 0
	for _, idw := range ideal {
		for _, gw := range got {
			if gw == idw {
				hitCount++
				break
			}
		}
	}
	return float64(hitCount) / float64(len(ideal))
}

func main() {
	// Example Corpus
	corpus := []string{
		"hello", "hell", "helicopter", "hero", "world",
		"how", "are", "you", "hello", "war", "hello",
	}

	// Metrics: Build tries and measure insertion time and memory
	startMem := getMemoryUsage()
	startTime := time.Now()

	trieA1 := NewTrieA1()
	for _, w := range corpus {
		trieA1.Insert(w)
	}
	trieA1.BuildBigramTable(corpus)

	buildTimeA1 := time.Since(startTime)
	endMemA1 := getMemoryUsage()
	memoryUsedA1 := endMemA1 - startMem

	// Algorithm 2 build
	startMem = getMemoryUsage()
	startTime = time.Now()

	trieA2 := initTriesA2()
	for _, w := range corpus {
		trieA2.Insert(w)
	}

	buildTimeA2 := time.Since(startTime)
	endMemA2 := getMemoryUsage()
	memoryUsedA2 := endMemA2 - startMem

	// Query metrics
	prefix := "he"
	k := 3

	// Algorithm_1 query
	startTime = time.Now()
	suggestionsA1 := trieA1.Autocomplete(prefix, k)
	queryTimeA1 := time.Since(startTime)

	var wordsA1 []string
	for _, s := range suggestionsA1 {
		wordsA1 = append(wordsA1, s.word)
	}

	// Algorithm_2 query
	startTime = time.Now()
	suggestionsA2 := trieA2.Autocomplete(prefix)
	if len(suggestionsA2) > k {
		suggestionsA2 = suggestionsA2[:k]
	}
	queryTimeA2 := time.Since(startTime)

	// For suggestion quality, define an ideal top-3 completions:
	// Let's say based on known frequency/context, we expect: ["hello", "helicopter", "hell"]
	ideal := []string{"hello", "helicopter", "hell"}

	qualityA1 := measureSuggestionQuality(wordsA1, ideal)
	qualityA2 := measureSuggestionQuality(suggestionsA2, ideal)

	// Print results
	fmt.Println("------ Algorithm 1 (Contextual) Metrics ------")
	fmt.Printf("Build Time: %v\n", buildTimeA1)
	fmt.Printf("Memory Used (bytes): %d\n", memoryUsedA1)
	fmt.Printf("Query Time: %v\n", queryTimeA1)
	fmt.Printf("Suggestions: %v\n", wordsA1)
	fmt.Printf("Suggestion Quality: %.2f\n\n", qualityA1)

	fmt.Println("------ Algorithm 2 (Frequency) Metrics ------")
	fmt.Printf("Build Time: %v\n", buildTimeA2)
	fmt.Printf("Memory Used (bytes): %d\n", memoryUsedA2)
	fmt.Printf("Query Time: %v\n", queryTimeA2)
	fmt.Printf("Suggestions: %v\n", suggestionsA2)
	fmt.Printf("Suggestion Quality: %.2f\n", qualityA2)
}
