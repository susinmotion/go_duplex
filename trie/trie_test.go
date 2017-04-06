package trie

import "testing"
import "fmt"

func TestCreate(*testing.T) {
	t := NewTrie()
	fmt.Printf("%b", (t == &Trie{}))
	//Output:
	//false
}
