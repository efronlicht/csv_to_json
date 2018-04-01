package main

import (
	"errors"
	"log"
	"os"
	"path/filepath"
)

//monitor the the input directory for new files ending in .csv
//for each file in the directory the
//create a file in outfolder writing output as JSON to
//the output directory
func main() {
	dirs, ok := getDirsFromArgs(os.Args)
	if !ok {
		log.Fatal("must specify exactly three directories: CSV_INPUT, JSON_OUTPUT, ERROR_OUPUT")
	}
	if err := dirs.jsonToCSV(); err != nil {
		log.Fatal(err)
	}
}

type dirs struct {
	in, out, err string
}

func (d dirs) jsonToCSV() error {
	csvs := make(chan string)
	go findNewCSVs(d.in, csvs)

	for csv := range csvs {
		if err := d.processCSV(csv); err != nil {
			return err
			// this is not especially fail-tolerant;
			// a problem interfacing with the OS on any given file
			// will kill the whole program.
		}
	}
	return nil
}

func (d dirs) processCSV(path string) error {
	records, parseErrs, err := readCSV(path)
	if err != nil {
		return err
	}

	outFile := d.outPath(path)
	if err := writeRecordsToFileAsJSON(outFile, records); err != nil {
		return err
	}

	errFile := d.errPath(path)
	if err := writeErrorsToFileAsCSV(errFile, parseErrs); err != nil {
		return err
	}

	return nil
}

/*deleteCSV deletes a file if and only if it's path ends in .csv
this isn't really necessary, as we guarantee our paths end in '.csv' otherwhere in the program,
but it never hurts to have a little extra protection when you're working with
the filesystem*/
func deleteCSV(path string) error {
	if filepath.Ext(path) != ".csv" {
		return errors.New("called deleteCSV on a filename that doens't end in .csv")
	}
	return os.Remove(path)
}

//get input, output, and error directories
func getDirsFromArgs(args []string) (dirs, bool) {
	if len(args) != 3 {
		return dirs{}, false
	}
	return dirs{args[0], args[1], args[2]}, true
}

func (d *dirs) errPath(filename string) string {
	return filepath.Join(d.err, filepath.Base(filename))
}
func (d *dirs) outPath(filename string) string {
	return filepath.Join(d.out, filepath.Base(filename))
}
func (d *dirs) inPath(filename string) string {
	return filepath.Join(d.in, filepath.Base(filename))
}
