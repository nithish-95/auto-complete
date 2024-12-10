package main

import (
	"testing"
	"time"
)

// Helper function to build Algorithm_1 trie and bigram table
func buildAlg1Trie(corpus []string) *TrieA1 {
	trie := NewTrieA1()
	for _, w := range corpus {
		trie.Insert(w)
	}
	trie.BuildBigramTable(corpus)
	return trie
}

// Helper to build Algorithm_2 trie
func buildAlg2Trie(corpus []string) *TriesA2 {
	trie := initTriesA2()
	for _, w := range corpus {
		trie.Insert(w)
	}
	return trie
}

// Test Case 1: Basic Prefix (No Context)
func TestBasicPrefixNoContext(t *testing.T) {
	corpus := []string{"hello", "hell", "helicopter", "hero", "world"}
	prefix := "he"
	ideal := []string{"hello", "helicopter", "hell"}

	trieA1 := buildAlg1Trie(corpus)
	trieA2 := buildAlg2Trie(corpus)

	suggestionsA1 := trieA1.Autocomplete(prefix, 5)
	suggestionsA2 := trieA2.Autocomplete(prefix)

	// Convert A1 suggestions to []string
	var wordsA1 []string
	for _, s := range suggestionsA1 {
		wordsA1 = append(wordsA1, s.word)
	}

	qA1 := measureSuggestionQuality(wordsA1, ideal)
	qA2 := measureSuggestionQuality(suggestionsA2, ideal)

	if len(wordsA1) == 0 || len(suggestionsA2) == 0 {
		t.Errorf("Expected non-empty suggestions for prefix '%s'", prefix)
	}

	// We only check that we have correct prefix matches and a reasonable quality score.
	if qA1 <= 0.0 {
		t.Errorf("Algorithm_1 did not match any ideal suggestions.")
	}
	if qA2 <= 0.0 {
		t.Errorf("Algorithm_2 did not match any ideal suggestions.")
	}
}

// Test Case 2: Contextual Ranking
func TestContextualRanking(t *testing.T) {
	corpus := []string{"hello", "hello", "hell", "helicopter", "hero", "world", "how", "are", "you", "hello", "war"}
	prefix := "he"
	// Suppose "hello" was the previous word, if "hell" often follows "hello" more than others,
	// Algorithm_1 should rank "hell" higher. Let's define ideal: "hell", "helicopter" as top due to context.
	// Without actual user context passed, we rely on bigram frequencies.
	ideal := []string{"hell", "helicopter"}

	trieA1 := buildAlg1Trie(corpus)
	trieA2 := buildAlg2Trie(corpus)

	suggestionsA1 := trieA1.Autocomplete(prefix, 5)
	suggestionsA2 := trieA2.Autocomplete(prefix)

	var wordsA1 []string
	for _, s := range suggestionsA1 {
		wordsA1 = append(wordsA1, s.word)
	}

	qA1 := measureSuggestionQuality(wordsA1, ideal)
	qA2 := measureSuggestionQuality(suggestionsA2, ideal)

	// If Algorithm_1 is contextual, qA1 should be >= qA2 in most cases.
	if qA1 < qA2 {
		t.Errorf("Algorithm_1 should provide better or equal contextual matches, got qA1: %f, qA2: %f", qA1, qA2)
	}
}

// Test Case 3: Non-Existent Prefix
func TestNonExistentPrefix(t *testing.T) {
	corpus := []string{"hello", "hell", "helicopter", "hero", "world"}
	prefix := "xyz"

	trieA1 := buildAlg1Trie(corpus)
	trieA2 := buildAlg2Trie(corpus)

	suggestionsA1 := trieA1.Autocomplete(prefix, 5)
	suggestionsA2 := trieA2.Autocomplete(prefix)

	if len(suggestionsA1) != 0 {
		t.Errorf("Expected empty result for Algorithm_1 with prefix '%s'", prefix)
	}
	if len(suggestionsA2) != 0 {
		t.Errorf("Expected empty result for Algorithm_2 with prefix '%s'", prefix)
	}
}

// Test Case 4: Long Prefix Single Match
func TestLongPrefixSingleMatch(t *testing.T) {
	corpus := []string{"helicopter", "hello", "hero"}
	prefix := "helico"

	trieA1 := buildAlg1Trie(corpus)
	trieA2 := buildAlg2Trie(corpus)

	suggestionsA1 := trieA1.Autocomplete(prefix, 5)
	suggestionsA2 := trieA2.Autocomplete(prefix)

	if len(suggestionsA1) == 0 || suggestionsA1[0].word != "helicopter" {
		t.Errorf("Algorithm_1 expected 'helicopter' for prefix '%s'", prefix)
	}
	if len(suggestionsA2) == 0 || suggestionsA2[0] != "helicopter" {
		t.Errorf("Algorithm_2 expected 'helicopter' for prefix '%s'", prefix)
	}
}

// Test Case 5: Large Corpus Performance Test
// This test is more about performance - run only if you want to measure.
// We'll just measure build times and ensure no errors occur.

func TestLargeCorpus(t *testing.T) {
	// Generate a large corpus (simple repetition)
	var corpus []string
	for i := 0; i < 10000; i++ {
		corpus = append(corpus, "he"+time.Now().Format("150405000000")) // unique words
	}
	prefix := "he"

	startTime := time.Now()
	trieA1 := buildAlg1Trie(corpus)
	buildA1Time := time.Since(startTime)

	startTime = time.Now()
	trieA2 := buildAlg2Trie(corpus)
	buildA2Time := time.Since(startTime)

	suggestionsA1 := trieA1.Autocomplete(prefix, 10)
	suggestionsA2 := trieA2.Autocomplete(prefix)

	// We don't have an ideal here, just checking no error and performance.
	if len(suggestionsA1) == 0 {
		t.Errorf("Algorithm_1 returned no suggestions. Possibly an error.")
	}
	if len(suggestionsA2) == 0 {
		t.Errorf("Algorithm_2 returned no suggestions. Possibly an error.")
	}

	// Just print performance info (not strictly pass/fail).
	t.Logf("Algorithm_1 Build Time: %v, Algorithm_2 Build Time: %v", buildA1Time, buildA2Time)
}
