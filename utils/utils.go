package utils

import (
	"fmt"
	"os"
	"time"
)

func PrintTime(start time.Time) {
	fmt.Println(time.Since(start))
}

func Checkerr(err error) {
	if err != nil {
		panic(err)
	}
}
func CheckConfig(err error, filename string) {
	if err != nil {
		fmt.Println("Config file", filename, "not found. Specify config file location using -C")
		os.Exit(1)
	} else {
		fmt.Sprintf("Parsed config file", filename)
	}
}

var curDir, err = os.Getwd()

func CheckOutput(filename string) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		os.Mkdir(filename, 0755)
	}
	if err != nil {
		fmt.Println("Cannot write to ", filename, ". Specify an output directory using -O")
		os.Exit(1)
	} else {
		if filename[0] == '/' {
			fmt.Println("Output can be found in folder", filename)
		} else {
			fmt.Println("Output can be found in", curDir+"/"+filename)

		}
	}
}

func ReverseComplement(line string) string {
	nucs := map[byte]byte{
		'A': 'T',
		'C': 'G',
		'G': 'C',
		'T': 'A',
		'N': 'N',
	}
	out := make([]byte, 0, len(line))
	for i := len(line) - 1; i >= 0; i-- {
		out = append(out, nucs[line[i]])
	}
	return string(out)
}
