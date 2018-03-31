package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

type dirs struct {
	in, out, err string
}

//monitor the the input directory for new files ending in .csv
//for each file in the directory the
//create a file in outfolder writing output as JSON to
//the output directory
func main() {
	in, out, errDir, ok := getArgs()
	if !ok {
		log.Fatal("must specify exactly three directories: CSV_INPUT, JSON_OUTPUT, ERROR_OUPUT")
	}

	filePaths := make(chan string)
	go findNewCSVs(in, filePaths)

	for csvFile := range filePaths {
		records, parseErrs, err := readCSV(csvFile)
		if err != nil {
			log.Fatal(err)
		}

		outFile := outPath(out, csvFile)
		if err := writeJSON(outFile, records); err != nil {
			log.Fatal(err)
		}

		errFile := errPath(errDir, csvFile)
		if err := writeErrorCSV(errFile, parseErrs); err != nil {
			log.Fatal(err)
		}

		if err := os.Remove(csvFile); err != nil {
			log.Fatal(err)
		}
	}
}

//get input, output, and error directories
func getArgs() (in, out, errDir string, ok bool) {
	if len(os.Args) != 3 {
		return
	}
	return os.Args[0], os.Args[1], os.Args[2], true
}

func errPath(dir, filename string) string {
	return filepath.Join(dir, filepath.Base(filename))
}

func outPath(dir, filename string) string {
	base := strings.TrimSuffix(filepath.Base(filename), ".csv")
	return filepath.Join(dir, base+".json") // we know the last 4 characters are EXACTLY .csv
}
