package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

func errorsToCSV(errs []csv.ParseError) []byte {
	lines := make([]string, len(errs))
	for i, err := range errs {
		lines[i] = fmt.Sprintf("%d,%s", err.Line, err.Err)
	}
	const header = "LINE_NUM,ERROR_MSG"
	body := strings.Join(lines, "\n")
	return []byte(header + "\n" + body)
}

func writeErrorCSV(filename string, errs []csv.ParseError) error {
	return ioutil.WriteFile(filename, errorsToCSV(errs), 0644)
}

func writeJSON(filename string, records []Record) error {
	asJSON, err := json.Marshal(records)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, asJSON, 0644)
}
