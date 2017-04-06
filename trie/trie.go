package trie

import "strings"
import "fmt"
import "sync"

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type LeafChunk struct {
	Lock    sync.Mutex
	Forward *LeafData
	Reverse *LeafData
	Super   *LeafData
}
type LeafData struct {
	Count    int
	Target   string
	Variants Variants
	Indel    Indel
	HasVar   bool
	HasIndel bool
	Trash    bool
}

func NewLeafChunk() *LeafChunk {
	l := &LeafChunk{
		Forward: NewLeafData(""),
		Reverse: NewLeafData(""),
		Super:   NewLeafData(""),
	}
	return l
}

type Variants map[int]VariantPos

func (l *LeafData) AddIndel(position int, length int) {
	if l.Indel.Length == 0 {
		l.HasIndel = true
		l.Indel = Indel{Pos: position, Length: length}
	} else if l.Indel.Length != length || l.Indel.Pos != position {
		l.Trash = true //we should fix this to allow for different indels
	}
}

func NewLeafData(data string) *LeafData {
	l := &LeafData{
		Target:   data,
		Variants: make(map[int]VariantPos),
		Indel:    Indel{},
		Trash:    false,
	}
	return l
}

func (v Variants) AddVariant(position int, target byte, actual byte) {
	key := position
	_, ok := v[key]
	if !ok {
		v[key] = NewVariantPos(NewBase(target))
	}
	b := NewBase(actual)
	_, ok = v[key].Counts[b]
	if !ok {
		v[key].Counts[b] = 0
	}
	v[key].Counts[b] = v[key].Counts[b] + 1
	if v[key].Target == b {
		fmt.Println("same bases", string(target), string(actual), Base2Str([]Base{v[key].Target}), Base2Str([]Base{b}))
	}

}

func (l *LeafData) Update(s string) {
	if l.Target == "" {
		l.Target = s
	} else {
		for i := 0; i < min(len(s), len(l.Target)); i++ {
			if s[i] != l.Target[i] {
				//fmt.Printf("addinvvariant")
				l.HasVar = true
				defer l.Variants.AddVariant(i, l.Target[i], s[i])
			}
		}
	}
	l.Count = l.Count + 1
}

type Node struct {
	lock     sync.Mutex
	value    Base
	Children [5]*Node
	Data     *LeafChunk
}

type Base int

func NewBase(c byte) Base {
	b := strings.IndexByte("ACGTN", c)
	return Base(b)
}

func Str2Base(s string) []Base {
	b := make([]Base, len(s))
	for i := 0; i < len(s); i++ {
		b[i] = NewBase(s[i])
	}
	return b
}

func Base2Str(b []Base) string {
	s := make([]byte, len(b))
	for i := 0; i < len(b); i++ {
		s[i] = "ACGTN"[b[i]]
	}
	return string(s)
}

type Trie struct {
	Root *Node
}

func NewTrie() *Trie {
	t := &Trie{
		Root: &Node{},
	}
	return t
}

func (t *Trie) Insert(key string, data string, strand string) *Node {
	n := t.Root
	for i := 0; i < len(key); i++ {
		n.lock.Lock()
		b := NewBase(key[i])
		if n.Children[b] == nil {
			n.Children[b] = &Node{value: b}
		}
		n.lock.Unlock()
		n = n.Children[b]
	}
	n.lock.Lock()
	if n.Data == nil {
		n.Data = NewLeafChunk()
	}
	if strand == "+" {
		n.Data.Forward.Update(data)
	} else if strand == "-" {
		n.Data.Reverse.Update(data)
	}
	n.lock.Unlock()
	return n
}

func (t *Trie) Traverse(n *Node, key []Base) {
	if n != nil {
		key = append(key, n.value)
		if n.Data != nil {
			fmt.Printf(Base2Str(key)[1:] + "\n")
			//fmt.Printf(Base2Str(n.Data.consensus))
			fmt.Println(n.Data)
			return
		} else {
			for i := 0; i < len(n.Children); i++ {
				t.Traverse(n.Children[i], key)
			}
		}

	}
}

func (curChunk *LeafChunk) Reconcile() {
	curChunk.Forward.UpdateConsensus(1) //I need to fix this
	curChunk.Reverse.UpdateConsensus(1)
	curChunk.Super.Target = curChunk.Forward.Target
	tmp := []byte(curChunk.Super.Target)
	mismatch := false
	for i := 0; i < len(curChunk.Forward.Target); i++ {
		if curChunk.Forward.Target[i] != curChunk.Reverse.Target[i] {
			//****record the variant between the strands
			tmp[i] = 'N'
			if (curChunk.Forward.Target[i] != 'N') && (curChunk.Reverse.Target[i] != 'N') {
				mismatch = true
			}
		}
	}
	if mismatch == true {
		fmt.Println(curChunk.Forward.Target, "-mismatch-", curChunk.Reverse.Target)
	}
	curChunk.Lock.Lock()
	curChunk.Super.Target = string(tmp) //also the targets should definitely be byte arrays. this is just annoying
	curChunk.Lock.Unlock()
}

func (curData *LeafData) UpdateConsensus(threshold int) {
	for pos, varpos := range curData.Variants {
		for observed, count := range varpos.Counts {
			if (float64(count) / float64(curData.Count)) >= float64(threshold) { //we have to deal with thresholding here. Let's say if it's not more than 100%
				temp := []byte(curData.Target)
				temp[pos] = "ACGTN"[observed]
				curData.Target = string(temp)
			} else {
				temp := []byte(curData.Target)
				temp[pos] = 'N'
				curData.Target = string(temp)
			}
		}
	}
}

func PPrintData(l *LeafData) {

}

type VariantPos struct {
	Target Base
	Counts map[Base]int
}

func NewVariantPos(target Base) VariantPos {
	v := VariantPos{
		Target: target,
		Counts: make(map[Base]int),
	}
	return v
}

type Indel struct {
	Pos    int
	Length int
}
