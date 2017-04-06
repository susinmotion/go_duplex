package parser

import (
	"bufio"
	"compress/gzip"
	"duplex/experiment"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	//"fmt"
)

func ReverseComplement(line string) string {
	nucs := map[byte]byte{
		'A': 'T',
		'C': 'G',
		'G': 'C',
		'T': 'A',
		'N': 'N',
		'a': 'T',
		'c': 'G',
		'g': 'C',
		't': 'T',
	}
	out := make([]byte, 0, len(line))
	for i := len(line) - 1; i >= 0; i-- {
		out = append(out, nucs[line[i]])
	}
	return string(out)
}
func Parse(wg *sync.WaitGroup, line string, experiments *map[string]*experiment.Experiment, direction string) {
	for _, v := range *experiments {
		fa := strings.Index(line, v.Metadata.FAlign)
		if fa == v.Metadata.BarcodeLen {
			name := v.Metadata.Name
			barcode := line[:fa]
			sequence := line[fa+len(v.Metadata.FAlign):]
			(*experiments)[name].AddSequence(barcode, sequence, direction)
			break
		}
	}
	wg.Done()
}

func Read(wg *sync.WaitGroup, filename string, experiments *map[string]*experiment.Experiment, direction string) {
	file, err := os.Open(filename)
	checkerr(err)
	fz, err := gzip.NewReader(file)
	checkerr(err)
	r := bufio.NewReader(fz)

	var wg2 sync.WaitGroup
	count := 0
	for {
		line, err := r.ReadString('\n')
		if err == io.EOF {
			break
		}
		checkerr(err)
		if count%4 == 1 {
			wg2.Add(1)
			line = strings.Trim(line, "\n ")
			go Parse(&wg2, line, experiments, direction)
		}
		count++
	}
	wg2.Wait()
	wg.Done()
}

func checkerr(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
