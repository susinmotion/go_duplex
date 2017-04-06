package main

import "duplex/experiment"
import "fmt"
import "duplex/config"
import "duplex/fastq"
import "time"
import "flag"
import "os"
import "path"
import "sync"
import "runtime/pprof"
import "strconv"
import "duplex/utils"
import "bowtie/bowtie"

var curDir, err = os.Getwd() //this is going to need to be os.Exec() in go 1.8
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var outputPtr = flag.String("O", curDir, "complete path to output directory")
var configPtr = flag.String("C", path.Join(curDir, "config.json"), "complete path to config file")

func main() {
	start := time.Now()
	defer utils.PrintTime(start)

	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		utils.Checkerr(err)
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	utils.CheckOutput(*outputPtr)
	//also check fasta and reference
	c, _ := config.ReadConfig(*configPtr)
	fmt.Println(c)

	es := experiment.Config2Exps(c)
	var wg2 sync.WaitGroup
	var wg sync.WaitGroup
	for i := 0; i < len(c.InputFiles); i++ {
		wg.Add(1)
		parser.Read(&wg, c.InputFiles[i], &es, "+")
	}
	for i := 0; i < len(c.RevInputFiles); i++ {
		wg.Add(1)
		parser.Read(&wg, c.RevInputFiles[i], &es, "-")
	}
	wg.Wait()

	for name, ex := range es {
		for _, thresh := range ex.Metadata.Thresholds {
			outfile := *outputPtr + "/" + name + "_" + strconv.Itoa(thresh)
			wg2.Add(1)
			go func(wg *sync.WaitGroup, outfile string, thresh int, ex *experiment.Experiment) {
				ex.Aggregate(outfile, thresh)
				fmt.Println("aggregated")
				bowtie.BowtieCheckVariants(outfile+"tmp", ex.Metadata.ReferenceGenome, ex.Metadata.SequenceLength, outfile)
				wg.Done()
			}(&wg2, outfile, thresh, ex)
		}
	}
	fmt.Println("parsed")
	wg2.Wait()
}

/*
func parseFiles(wg *sync.WaitGroup, filename string, strand string, es map[string]*experiment.Experiment){
    fmt.Println("parsing",filename)
    parser.Read(filename, &es, strand)
    wg.Done()
}*/
