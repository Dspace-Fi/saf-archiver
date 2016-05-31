package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const contentfile = "contents"

func is_leaf(path string) bool {

	fs, err := ioutil.ReadDir(path)
	if err != nil {
		return false // TODO: errors may be masked here
	}

	for _, f := range fs {
		if f.IsDir() {
			return false
		}
	}
	return true
}
func list_leaf_directories(base string) []string {

	var paths []string

	var walker = func(path string, _ os.FileInfo, _ error) error {
		if is_leaf(path) {
			paths = append(paths, path)
		}
		return nil
	}

	filepath.Walk(base, walker)

	return paths
}

func cp(src string, dest string) error {

	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()

	d, err := os.Create(dest)
	if err != nil {
		return err
	}

	if _, err := io.Copy(d, s); err != nil {
		d.Close()
		return err
	}
	return d.Close()
}

func add_file_metadata(fn string, md string) error {
	f, err := os.OpenFile(fn, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	f.WriteString(md + "\n")

	return nil
}

func add_file(fn string, dir string, mds []string) {
	fmt.Println(fn, dir)
	bn := filepath.Base(fn)
	target := filepath.Join(dir, bn)
	if _, err := os.Stat(target); err == nil {
		fmt.Fprintf(os.Stderr, "File %v exists already, skipping...\n", target)
		return
	}

	if err := cp(fn, target); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot copy file %v to %v: %v\n", fn, target, err)
		return
	}

	md := bn

	if len(mds) != 0 {
		s := strings.Join(mds, "\t")
		md = md + "\t" + s
	}

	mdfn := filepath.Join(dir, contentfile) // TODO declare in somewhere more visible place
	if err := add_file_metadata(mdfn, md); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing `contents` file: %v\n", err)
		// TODO: perhaps delete copied file?
		return
	}
}

func main() {

	metadata := flag.String("m", "", "Additional metadata for file, optionally separated with comma")

	flag.Parse()
	args := flag.Args()

	if len(args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %v [options] added-file target-directory\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
		os.Exit(1)
	}

	fn := args[0]
	basedir := args[1]

	var mds []string

	if *metadata != "" {
		mds = strings.Split(*metadata, ",")
	} else {
		mds = nil
	}

	fmt.Println(mds)

	dirs := list_leaf_directories(basedir)
	for _, d := range dirs {
		add_file(fn, d, mds)
	}
}
