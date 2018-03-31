package main

import (
	"encoding/csv"
	"io"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/eyecuelab/kit/set" // I wrote this StringSet type
)

//coroutine that sends newly discovered filenames ending in .csv to out channel
func findNewCSVs(dir string, out chan<- string) {
	seen := make(set.String)
	for ; ; time.Sleep(1 * time.Second) {
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			panic(err) // if we can't read the directory, this isn't recoverable
		}
		for _, file := range files {
			name := file.Name()
			if path.Ext(name) == ".csv" && !seen.Contains(name) {
				seen.Add(name)
				out <- dir + name
			}
		}
	}
}

//read a CSV, extracting all valid records (for valid lines) and parseErrors (for invalid lines)
func readCSV(filepath string) (records []Record, parseErrs []csv.ParseError, err error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()
	var r = csv.NewReader(f)
	_, _ = r.Read() // skip header
	line := 1
	for fields, err := r.Read(); err != io.EOF; fields, err = r.Read() {
		line++ // the header is line 1; our first real line is line 2
		if record, err := toRecord(fields); err != nil {
			parseErrs = append(parseErrs, csv.ParseError{Line: line, Err: err})
		} else {
			records = append(records, record)
		}
	}
	return records, parseErrs, nil
}
