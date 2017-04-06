package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mds/utils"
	"os"
	"path/filepath"
)

type Metadata struct {
	BarcodeLen      int    `json:"barcode length"`
	Name            string `json:"gene"`
	FAlign          string `json:"forward align"`
	Thresholds      []int  `json:"thresholds"`
	RCFAlign        string
	ReferenceGenome string `json:"reference genome"`
	SequenceLength  int    `json:"sequence length"`
}

type Config struct {
	InputFiles    []string   `json:"input files"`
	Mode          string     `json:"mode"`
	RevInputFiles []string   `json:"input files reverse"`
	Genes         []Metadata `json:"genes"`
}
type ConfigParser func(data []byte) Config

func ParseJSON(data []byte) Config {
	var res Config
	err := json.Unmarshal(data, &res)
	utils.Checkerr(err)
	//fmt.Printf("PARSED\t%d records\n", len(res))
	//fmt.Printf(string(data))
	return res
}

func ParseYAML(data []byte) Config {
	return Config{}
}

func ReadConfig(filename string) (Config, error) {
	var parser ConfigParser = ParseJSON
	data, err := ioutil.ReadFile(filename)
	utils.CheckConfig(err, filename)
	ext := filepath.Ext(filename)
	switch ext {
	case ".json":
		parser = ParseJSON
	case ".yml", ".yaml":
		parser = ParseYAML
	default:
		return Config{}, errors.New("ERROR: unknown config file format. Use .json or .yaml")
	}
	return parser(data), nil
}

func PrintConfig(c Config) {
	bytes, err := json.Marshal(&c)
	if err != nil {
		fmt.Println("ERROR converting Metadata to .json")
		os.Exit(1)
	}
	fmt.Println(string(bytes))
	//fmt.Println(c.Genes)
}
