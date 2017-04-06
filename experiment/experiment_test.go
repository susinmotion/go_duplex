package experiment

//import "testing"
import "mds/trie"
import "fmt"

func ExampleInsert() {
	fmt.Printf("elloo")
	e := Experiment{Data: trie.NewTrie(), Metadata: Metadata{MinThresh: 3, Target: "AACG"}, c: make(chan *trie.LeafData)}
	fmt.Printf("hello")
	e.addSequence("CAAAAC", "NACG")
	e.addSequence("CAAAAC", "NACG")
	e.addSequence("CAAAAC", "NACG")
	e.Data.Traverse(e.Data.Root, nil)
	// Output:
	// CAAAAC
	//&{3 map[0:map[4:3]] {0 0} false}
}
