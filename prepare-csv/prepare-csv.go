package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type config struct {
	InputSeparator  string   `json:"input-separator"`
	OutputSeparator string   `json:"output-separator"`
	SplitSeparator  string   `json:"split-separator"`
	Columns         []Column `json:"columns"`
}

type Column struct {
	From    int    `json:"from"`
	To      int    `json:"to"`
	Title   string `json:"title"`
	SplitBy string `json:"split-by"`
}

func printSlice(xs []string, sep string, out io.Writer) {
	l := len(xs)
	for i, e := range xs {
		fmt.Fprint(out, e)
		if (i + 1) < l {
			fmt.Print(sep)
		}
	}
	fmt.Fprintln(out)
}

func makeHeader(cols []Column) []string {
	target := make([]string, len(cols))

	for _, e := range cols {
		target[e.To] = e.Title
	}
	return target
}

func processRecord(record []string, cols []Column, splitter string, sep string) []string {

	target := make([]string, len(cols))
	for _, e := range cols {
		s := record[e.From]

		// replace splitter strings, if necessary
		if e.SplitBy != "" {
			s = strings.Replace(s, e.SplitBy, splitter, -1)
		}
		// escape string, if it contains separator
		if strings.Contains(s, sep) {
			s = "\"" + s + "\""
		}

		target[e.To] = s
	}
	return target
}

func main() {

	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %v config-file input-file\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	// read config
	var conf config
	f, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot read config file: %v\n", err)
		os.Exit(1)
	}
	err = json.Unmarshal(f, &conf)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot parse config file: %v\n", err)
		os.Exit(1)
	}

	// set defaults, if not specified in configuration
	if conf.InputSeparator == "" {
		conf.InputSeparator = ";"
	}
	if conf.OutputSeparator == "" {
		conf.OutputSeparator = ";"
	}
	if conf.SplitSeparator == "" {
		conf.SplitSeparator = ";"
	}

	// process file
	fn := os.Args[2]

	infile, err := os.Open(fn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot open input file: %v\n", err)
		os.Exit(1)
	}
	defer infile.Close()

	r := csv.NewReader(infile)
	r.Comma = rune(conf.InputSeparator[0]) // TODO works only with 8-bits
	r.LazyQuotes = true

	records, err := r.ReadAll()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot read records: %v\n", err)
		os.Exit(1)
	}

	printSlice(makeHeader(conf.Columns), conf.OutputSeparator, os.Stdout)

	// process records
	for _, e := range records {
		printSlice(processRecord(e, conf.Columns, conf.SplitSeparator, conf.OutputSeparator),
			conf.OutputSeparator, os.Stdout)
	}
}
