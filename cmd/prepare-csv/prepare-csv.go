package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/dspace-fi/saf-archiver"
)

type config struct {
	InputSeparator  string   `json:"input-separator"`
	OutputSeparator string   `json:"output-separator"`
	SplitSeparator  string   `json:"split-separator"`
	Columns         []Column `json:"columns"`
}

type Column struct {
	From    int      `json:"from"`
	Discard bool     `json:"discard"`
	Title   string   `json:"title"`
	SplitBy string   `json:"split-by"`
	Filters []string `json:"filters"`
}

func makeHeader(cols []Column) []string {
	var target []string

	for _, e := range cols {
		if !e.Discard {
			target = append(target, e.Title)
		}
	}
	return target
}

func processRecord(record []string, cols []Column, splitter string) []string {

	var target []string
	for _, e := range cols {

		if e.Discard {
			continue
		}

		s := record[e.From]

		// replace splitter strings, if necessary
		if e.SplitBy != "" {
			s = strings.Replace(s, e.SplitBy, splitter, -1)
		}

		// apply filters, if any
		if len(e.Filters) != 0 {
			for _, f := range e.Filters {

				fn, ok := filter.Filters[f]

				if ok {
					s = fn(s)
				} else {
					fmt.Fprintf(os.Stderr, "Cannot find filter %v! Aborting!\n", f)
					fmt.Fprintf(os.Stderr, "Available filters: ")
					for k := range filter.Filters {
						fmt.Fprintf(os.Stderr, "%v ", k)
					}
					fmt.Fprintln(os.Stderr)
					os.Exit(1)
				}
			}
		}

		target = append(target, s)

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
	r.Comma = rune(conf.InputSeparator[0]) // TODO works only with 8-bit characters
	r.LazyQuotes = true

	records, err := r.ReadAll()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot read records: %v\n", err)
		os.Exit(1)
	}

	w := csv.NewWriter(os.Stdout)
	w.Comma = rune(conf.OutputSeparator[0]) // TODO works only with 8-bit characters

	if err := w.Write(makeHeader(conf.Columns)); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot write header to CSV: %v\n", err)
		os.Exit(1)
	}

	// process records
	for _, e := range records {
		if err := w.Write(processRecord(e, conf.Columns, conf.SplitSeparator)); err != nil {
			fmt.Fprintf(os.Stderr, "Cannot write record to CSV: %v\n", err)
			os.Exit(1)
		}
	}

	w.Flush()
}
