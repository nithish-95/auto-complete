package main

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/nsf/termbox-go"
)

// Node struct represents each node in the trie with frequency
type Node struct {
	children    map[rune]*Node
	isEndOfWord bool
	frequency   int
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
			frequency:   0,
		},
	}
}

// Insert will take in a word and insert it into the trie, updating frequency
func (t *Tries) Insert(word string) {
	current := t.root
	for _, char := range word {
		node, ok := current.children[char]
		if !ok {
			node = &Node{
				isEndOfWord: false,
				children:    make(map[rune]*Node),
				frequency:   0,
			}
			current.children[char] = node
		}
		current = node
	}
	current.isEndOfWord = true
	current.frequency++ // Increment frequency only at the end of the word
}

// Get frequency of a word in the Trie
func (t *Tries) getFrequency(word string) int {
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

// Autocomplete returns a list of words that match the prefix, sorted by frequency
func (t *Tries) Autocomplete(prefix string) []string {
	current := t.root
	for _, char := range prefix {
		node, ok := current.children[char]
		if !ok {
			return []string{}
		}
		current = node
	}

	var results []string
	collectWords(current, prefix, &results)

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

// Collect words from a given node
func collectWords(node *Node, prefix string, results *[]string) {
	if node.isEndOfWord {
		*results = append(*results, prefix)
	}
	for char, child := range node.children {
		collectWords(child, prefix+string(char), results)
	}
}

func main() {
	// Initialize termbox
	err := termbox.Init()
	if err != nil {
		log.Fatalf("Failed to initialize termbox: %v", err)
	}
	defer termbox.Close()

	// Initialize the trie
	trie := initTries()
	words := []string{"hello", "hell", "helicopter", "hero", "world", "how", "are", "you"}
	for _, word := range words {
		trie.Insert(word)
	}

	var fullSentence string
	var currentWord string
	var suggestions []string

	render := func() {
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

		// Display the full sentence
		for i, char := range fullSentence {
			termbox.SetCell(i, 0, char, termbox.ColorDefault, termbox.ColorDefault)
		}

		// Display the current word being typed
		for i, char := range currentWord {
			termbox.SetCell(len(fullSentence)+i+1, 0, char, termbox.ColorGreen, termbox.ColorDefault)
		}

		// Display suggestions
		for i, suggestion := range suggestions {
			line := fmt.Sprintf("%d. %s", i+1, suggestion)
			for j, char := range line {
				termbox.SetCell(j, i+1, char, termbox.ColorDefault, termbox.ColorDefault)
			}
		}

		termbox.Flush()
	}

	for {
		render()

		ev := termbox.PollEvent()
		if ev.Type == termbox.EventKey {
			switch ev.Key {
			case termbox.KeyEsc:
				return // Exit on Escape key
			case termbox.KeyBackspace, termbox.KeyBackspace2:
				// Handle backspace
				if len(currentWord) > 0 {
					currentWord = currentWord[:len(currentWord)-1]
				} else if len(fullSentence) > 0 {
					words := strings.Fields(fullSentence)
					if len(words) > 0 {
						currentWord = words[len(words)-1]
						fullSentence = strings.Join(words[:len(words)-1], " ")
					}
				}
			case termbox.KeySpace:
				// Finalize the current word
				if currentWord != "" {
					fullSentence += currentWord + " "
					if trie.getFrequency(currentWord) == 0 {
						trie.Insert(currentWord)
					}
					currentWord = ""
				}
			case termbox.KeyTab:
				// Autocomplete the current word with the top suggestion
				if len(suggestions) > 0 {
					currentWord = suggestions[0]
				}
			case termbox.KeyEnter:
				// Finalize the sentence
				if currentWord != "" {
					fullSentence += currentWord + " "
					if trie.getFrequency(currentWord) == 0 {
						trie.Insert(currentWord)
					}
					currentWord = ""
				}
				fmt.Println("\nYour sentence:", fullSentence)
				return
			default:
				if ev.Ch != 0 {
					currentWord += string(ev.Ch)
				}
			}

			// Update suggestions for the current word
			suggestions = trie.Autocomplete(currentWord)
		}
	}
}
