package main

import (
	"bytes"
	"encoding/csv"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// TODO: override flag.Usage

type DublinCore struct {
	XMLName  xml.Name  `xml:"dublin_core"`
	Schema   string    `xml:"schema,attr,omitempty"`
	DCValues []DCValue `xml:"dcvalue"`
}

type DCValue struct {
	Element   string `xml:"element,attr"`
	Qualifier string `xml:"qualifier,attr,omitempty"`
	Language  string `xml:"language,attr,omitempty"`
	Value     string `xml:",innerxml"`
}

func makeDublinCore(schema string) *DublinCore {

	if schema == "dc" {
		return &DublinCore{}
	} else {
		return &DublinCore{Schema: schema}
	}
}

func makeDCValues(header string, value string) []DCValue {

	// TODO process header only once
	if value == "" {
		return nil
	}
	lang := ""

	// check for language
	xs := strings.Split(header, ":")

	if len(xs) > 1 {
		if len(xs) > 2 {
			fmt.Fprintf(os.Stderr, "Invalid header: %v can contain at most single language indicator ':'.\n",
				header)
			return nil
		}
		header = xs[0]
		lang = xs[1]
	}

	ys := strings.Split(header, ".")

	if len(ys) < 2 || len(ys) > 3 {
		fmt.Fprintf(os.Stderr, "Invalid header: %v, has %d elements.\n", header, len(ys))
		return nil
	}

	var dcvalues []DCValue

	vs := strings.Split(value, "||") // TODO from configuration?

	for _, v := range vs {

		dcvalue := DCValue{Element: ys[1], Value: v}

		if len(ys) == 3 {
			dcvalue.Qualifier = ys[2]
		} else {
			dcvalue.Qualifier = "none"
		}
		if lang != "" {
			dcvalue.Language = lang
		}
		dcvalues = append(dcvalues, dcvalue)
	}

	return dcvalues
}

func xmlFilename(s string) string {
	if s == "dc" {
		return "dublin_core.xml"
	} else {
		return "metadata_" + s + ".xml"
	}
}

func escapeEntities(s string) string {
	buf := new(bytes.Buffer)
	if err := xml.EscapeText(buf, []byte(s)); err == nil {
		return buf.String()
	} else {
		return ""
	}
}

func processRecord(xs []string, headers []string) map[string]*DublinCore {

	xmls := make(map[string]*DublinCore)

	for i, header := range headers {

		schema := strings.Split(header, ".")[0] // TODO check for error
		value := xs[i]

		if _, ok := xmls[schema]; ok {
			xmls[schema].DCValues = append(xmls[schema].DCValues,
				makeDCValues(header, escapeEntities(value))...)

		} else {
			xmls[schema] = makeDublinCore(schema)
			xmls[schema].DCValues = append(xmls[schema].DCValues,
				makeDCValues(header, escapeEntities(value))...)
		}
	}
	return xmls
}

func writeDC(dc *DublinCore, w io.Writer) {

	out, err := xml.MarshalIndent(dc, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating XML: %v\n", err)
	} else {
		w.Write([]byte(xml.Header))
		w.Write(out)
	}
}

func createDirectoryOrDie(dir string) {

	if _, err := os.Stat(dir); err == nil {
		fmt.Fprintf(os.Stderr, "Output directory '%v' exists already!\n", dir)
		os.Exit(1)
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot create directory: %v\n", err)
		os.Exit(1)
	}
}

func main() {

	trailfile := flag.String("t", "", "Filename where to output relative directories for each item")

	flag.Parse()

	args := flag.Args()

	if len(args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %v [options] input-file.csv archive_directory\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
		os.Exit(1)
	}

	// read records from input file
	infile, err := os.Open(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannon open input file: %v\n", err)
		os.Exit(1)
	}
	defer infile.Close()

	r := csv.NewReader(infile)
	r.Comma = ';' // TODO from config
	r.LazyQuotes = true

	records, err := r.ReadAll()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot read records: %v\n", err)
		os.Exit(1)
	}

	// create output directory
	basedir := args[1]
	createDirectoryOrDie(basedir)

	// create trail file
	var tf *os.File

	if *trailfile != "" {
		tf, err = os.Create(*trailfile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot create trail file %v (%v)\n",
				*trailfile, err)
			os.Exit(1)
		}
		defer tf.Close()
	}

	// process records
	headers := records[0]
	rest := records[1:]

	// sanity check for headers
	if len(headers) == 0 {
		fmt.Fprintf(os.Stderr, "No header fields: %v! Exiting.\n", headers)
		os.Exit(1)
	}

	for i, fields := range rest {

		// sanity check for record fields
		if len(fields) != len(headers) {
			fmt.Fprintf(os.Stderr,
				"Line %d, got only %d elements, expected %d! Skipping.\n",
				(i + 1), len(fields), len(headers))
			continue
		}

		// process a record
		xmls := processRecord(fields, headers)

		// create files
		dir := path.Join(basedir, "item_"+fmt.Sprintf("%03d", i))
		createDirectoryOrDie(dir)

		for k, v := range xmls {
			fn := xmlFilename(k)
			fn = path.Join(dir, fn)

			f, err := os.Create(fn)
			if err != nil {
				fmt.Fprintf(os.Stderr,
					"Cannot create file %v (%v)! Skipping.\n",
					fn, err)
				continue
			}
			writeDC(v, f)
			// update trailfile
			tf.WriteString(dir + "\n")
			f.Close()
		}
	}
}
