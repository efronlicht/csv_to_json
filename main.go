package main

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	DEFAULT_TIMEOUT = time.Hour
	TEST_TIMEOUT    = 250 * time.Millisecond
)

var is_test = false

func main() {
	var timeout time.Duration
	if is_test {
		timeout = TEST_TIMEOUT
	} else {
		timeout = DEFAULT_TIMEOUT
	}
	dirs, err := getDirsFromArgs(os.Args)
	if err != nil {
		log.Fatal(err)
	}
	if err := dirs.jsonToCSV(timeout); err != nil {
		log.Fatal(err)
	}
}

type dirs struct {
	in, out, err string
}

//monitor the the input directory for new files ending in .csv.
//for each file in the directory, create a file in outfolder writing output as JSON to
//the output directory
func (d dirs) jsonToCSV(timeout time.Duration) error {
	csvs := make(chan string)
	go findNewCSVs(d.in, csvs, timeout)

	for csv := range csvs { // this is an indefinite loop unless processCSV hits an error
		log.Print(csv)
		if err := d.processCSV(csv); err != nil {
			return err
			// this is not especially fail-tolerant;
			// a problem interfacing with the OS on any given file
			// will kill the whole program.
		}
	}
	return nil
}

//processCSV processes the indivudal CSV specified by path
//given directories $in, $out, and $err, and filename "foo.csv",
//read from "$in/foo.csv"
func (d dirs) processCSV(path string) error {
	records, parseErrs, err := readCSVFile(d.inPath(path))
	if err != nil {
		return err
	}
	log.Printf("read data from %s", d.inPath(path))

	if err = writeRecordsToFileAsJSON(d.outPath(path), records); err != nil {
		return err
	}
	log.Printf("wrote records as json to %s", d.outPath(path))
	if err = writeErrorsToFileAsCSV(d.errPath(path), parseErrs); err != nil {
		return err
	}
	if err := deleteCSV(d.inPath(path)); err != nil {
		return err
	}
	log.Printf("wrote parseErrors to %s", d.errPath(path))
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
func getDirsFromArgs(args []string) (d dirs, err error) {
	if len(args) != 4 { //args[0] is always the program
		return d, errors.New("must specify exactly three directories: CSV_INPUT, JSON_OUTPUT, ERROR_OUPUT")
	}
	var in, out, errDir string
	if in, err = filepath.Abs(args[1]); err != nil {
		return d, err
	}
	if out, err = filepath.Abs(args[2]); err != nil {
		return d, err
	}
	if errDir, err = filepath.Abs(args[3]); err != nil {
		return d, err
	}
	return dirs{in, out, errDir}, nil
}

func (d *dirs) errPath(filename string) string {
	return filepath.Join(d.err, filepath.Base(filename))
}
func (d *dirs) outPath(filename string) string {
	filename = strings.TrimSuffix(filename, ".csv") + ".json"
	return filepath.Join(d.out, filepath.Base(filename))
}
func (d *dirs) inPath(filename string) string {
	return filepath.Join(d.in, filepath.Base(filename))
}
