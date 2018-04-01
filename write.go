package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

const ERROR_HEADER = "LINE_NUM,ERROR_MSG"

//errorsToCSV formats errors as CSV in the following format:
//"LINE_NUM,ERROR_MSG"
//"2,cosmic ray"
func errorsToCSV(errs []csv.ParseError) []byte {
	lines := make([]string, len(errs)+1)
	lines[0] = ERROR_HEADER
	for i, err := range errs {
		lines[i+1] = fmt.Sprintf("%d,%s", err.Line, err.Err)
	}
	return []byte(strings.Join(lines, "\n"))
}

func writeErrorsToFileAsCSV(filename string, errs []csv.ParseError) error {
	return ioutil.WriteFile(filename, errorsToCSV(errs), 0644)
}

//writeRecordsToFileAsJSON marshalls the records into JSON and writes them to the file
//specified by filename, creating or truncating the file if necessary.
func writeRecordsToFileAsJSON(filename string, records []Record) error {
	asJSON, err := json.Marshal(records)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, asJSON, 0644)

}
