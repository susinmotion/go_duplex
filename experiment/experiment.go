package experiment

import (
	"duplex/config"
	"duplex/trie"
	"duplex/utils"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	//"os/exec"
	//"io/ioutil"
	//"bytes"
	//"time"
	//"duplex/utils"
)

type LeafDataArray []*trie.LeafChunk

type Experiment struct {
	Metadata       config.Metadata
	Data           *trie.Trie
	ImportantNodes LeafDataArray
	INLock         sync.Mutex
	SuperVariants  map[int]trie.VariantPos
	SuperIndels    map[int]trie.Indel
	DifVariants    map[int]trie.VariantPos
	Coverage       map[int]int
	Shifts         [4][5]int
	DataLock       sync.Mutex
}

func NewExperiment(metadata config.Metadata) *Experiment {
	e := &Experiment{
		Metadata:       metadata,
		Data:           trie.NewTrie(),
		ImportantNodes: LeafDataArray{},
		SuperVariants:  make(map[int]trie.VariantPos),
		SuperIndels:    make(map[int]trie.Indel),
		DifVariants:    make(map[int]trie.VariantPos),
		Coverage:       make(map[int]int),
		Shifts:         [4][5]int{},
	}
	return e
}

func Config2Exps(c config.Config) map[string]*Experiment {
	es := make(map[string]*Experiment)
	for i := 0; i < len(c.Genes); i++ {
		e := NewExperiment(c.Genes[i])
		es[c.Genes[i].Name] = e
	}
	return es
}

func (e *Experiment) AddNode(n *trie.LeafChunk) {
	e.INLock.Lock()
	e.ImportantNodes = append(e.ImportantNodes, n)
	e.INLock.Unlock()
}

func (e *Experiment) AddSequence(barcode string, sequence string, strand string) {
	n := e.Data.Insert(barcode, sequence, strand)
	if n.Data.Forward.Count+n.Data.Reverse.Count == e.Metadata.Thresholds[0] {
		e.AddNode(n.Data)
	}
}

func (e *Experiment) Aggregate(outfile string, thresh int) {

	//for each of fwd and reverse, we want to update the consensus based on the threshold.
	var count int
	var wg2 sync.WaitGroup
	//jobs := make(chan int,1000)
	ch := make(chan string)
	fo, err := os.Create(outfile + "tmp")
	utils.Checkerr(err)
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()

	wg2.Add(1)
	go func(c chan string, fo *os.File) {
		for data := range ch {
			fo.Write([]byte(data + "\n"))
		}
		wg2.Done()
	}(ch, fo)

	for i := 0; i < len(e.ImportantNodes); i++ {
		curChunk := e.ImportantNodes[i]
		if (curChunk.Forward.Count >= thresh) && (curChunk.Reverse.Count >= thresh) {
			count++
			//jobs<-1
			curChunk.Reconcile()
			ch <- curChunk.Super.Target
			/*pos, variants :=bowtieCheckVariants(curChunk.Super.Target, e.Metadata.ReferenceGenome) //this is the slowest part is there another way we can do this? yes.
			e.DataLock.Lock()
			_, ok :=e.Coverage[pos]
			if !ok{
				e.Coverage[pos]=0
			}
			e.Coverage[pos]++
			for pos, varpos:= range(variants){
				_, ok = e.SuperVariants[pos]
				if !ok{
					e.SuperVariants[pos]=trie.NewVariantPos(varpos.Target)
				}
				for obs, _:= range(variants[pos].Counts){
					e.Shifts[varpos.Target][obs]+=1
					e.SuperVariants[pos].Counts[obs]+=1
				}
			}
			e.DataLock.Unlock()
			<-jobs*/
			//write the thing to a file
		}

	}
	close(ch)
	wg2.Wait()
	//ppShifts(e.Shifts, outfile)
	//ppVariants(e.SuperVariants, outfile, count)
	//ppCoverage(e.Coverage, outfile)
}

func PpVariants(variants map[int]trie.VariantPos, outfile string, count int) {
	//pos base count
	var keys []int
	for k, _ := range variants {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	fo, err := os.Create(outfile)
	fmt.Println(err)
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()
	fo.Write([]byte(strconv.Itoa(count) + " important nodes\n"))
	for _, k := range keys {
		for kk, vv := range variants[k].Counts {
			fo.Write([]byte(strconv.Itoa(k) + " " + string("ACGTN"[variants[k].Target]) + "->" + string("ACGTN"[kk]) + " " + strconv.Itoa(vv) + "\n"))
		}
	}
}

func ppShifts(shifts [4][5]int, outfile string) {
	file, _ := os.Create(outfile + ".csv")
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()
	for _, value := range shifts {
		file.Write([]byte(strings.Trim(strings.Join(strings.Fields(fmt.Sprint(value)), ","), "[]") + "\n"))
	}

}

func PpCoverage(coverage map[int]int, outfile string) {
	var keys []int
	for k, _ := range coverage {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	fo, _ := os.Create(outfile + "_coverage.csv")

	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()
	for _, k := range keys {
		fo.Write([]byte(strconv.Itoa(k) + "," + strconv.Itoa(coverage[k]) + "\n"))
	}

}
