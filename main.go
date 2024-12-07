package main

import (
	"bufio"
	"fmt"
	"os"
)

// Node struct represents each node in the trie
type Node struct {
	children    map[rune]*Node
	isEndOfWord bool
}

// Tries struct represents the trie with a pointer to the root node
type Tries struct {
	root *Node
}

// Initialize trie using initTries function
func initTries() *Tries {
	return &Tries{
		root: &Node{
			isEndOfWord: false,
			children:    make(map[rune]*Node),
		},
	}
}

// Insert will take in a word and insert it into the trie
func (t *Tries) Insert(word string) {
	current := t.root
	for _, char := range word {
		node, ok := current.children[char]
		if !ok {
			node = &Node{
				isEndOfWord: false,
				children:    make(map[rune]*Node),
			}
			current.children[char] = node
		}
		current = node
	}
	current.isEndOfWord = true
}

// Search will search for a word in the trie
func (t *Tries) Search(word string) bool {
	current := t.root
	for _, char := range word {
		node, ok := current.children[char]
		if !ok {
			return false
		}
		current = node
	}
	return current.isEndOfWord
}

// Autocomplete will return a list of words that match the prefix
func (t *Tries) Autocomplete(prefix string) []string {
	current := t.root
	for _, char := range prefix {
		node, ok := current.children[char]
		if !ok {
			return []string{} // Prefix not found
		}
		current = node
	}

	// Recursively collect all words starting from the current node
	var results []string
	t.collectWords(current, prefix, &results)
	return results
}

// Helper function to collect words from a given node
func (t *Tries) collectWords(node *Node, prefix string, results *[]string) {
	if node.isEndOfWord {
		*results = append(*results, prefix)
	}
	for char, child := range node.children {
		t.collectWords(child, prefix+string(char), results)
	}
}

func main() {
	testTries := initTries()
	scanner := bufio.NewScanner(os.Stdin)

	// Insert sample words
	testTries.Insert("hello")
	testTries.Insert("hell")
	testTries.Insert("helicopter")
	testTries.Insert("hero")
	testTries.Insert("world")
	testTries.Insert("how")
	testTries.Insert("are")
	testTries.Insert("you")

	fmt.Println("Enter a prefix to autocomplete:")
	for scanner.Scan() {
		prefix := scanner.Text()
		autocompleteResults := testTries.Autocomplete(prefix)
		if len(autocompleteResults) > 0 {
			fmt.Printf("Autocomplete results for '%s': %v\n", prefix, autocompleteResults)
		} else {
			fmt.Printf("No words found for prefix '%s'\n", prefix)
		}
		fmt.Println("Enter a prefix to autocomplete:")
	}
}
