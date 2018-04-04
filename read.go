package main

import (
	"encoding/csv"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

	"github.com/eyecuelab/kit/set" // I wrote this library
)

//findNewCSVS is a coroutine that periodically scans a directory for new files ending in .csv,
//sending unique filenames to the out channel.
//this function does not return and calls log.Fatal on filesystem errors
func findNewCSVs(dir string, out chan<- string, timeout time.Duration) {
	seen := make(set.String)
	defer close(out)
	for start := time.Now(); time.Now().Sub(start) < timeout; time.Sleep(time.Second) {
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			log.Fatal(err)
		}
		for _, file := range files {
			name := file.Name()
			if path.Ext(name) == ".csv" && !seen.Contains(name) {
				seen.Add(name)
				out <- name
			}
		}
	}
}

//readCSV takes a path to a CSV, opens it,
//extracting all valid records (for valid lines) and parseErrors (for invalid lines)
func readCSVFile(filepath string) (records []Record, parseErrs []csv.ParseError, err error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)

	_, _ = r.Read()
	// we assume a header exists & skip it. note that
	// this will cause us to silently skip the first line of a headless but
	// otherwise properly formatted .csv

	line := 1
	for fields, err := r.Read(); err != io.EOF; fields, err = r.Read() {
		line++ //the header is line 1; our first real line is line 2
		if record, err := toRecord(fields); err != nil {
			parseErrs = append(parseErrs, csv.ParseError{Line: line, Err: err})
		} else {
			records = append(records, record)
		}
	}
	return records, parseErrs, nil
}
