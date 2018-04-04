package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

//integration test. sets Os.Args manually and runs main
func Test_main(t *testing.T) {
	const testCSVData = "INTERNAL_ID,FIRST_NAME,MIDDLE_NAME,LAST_NAME,PHONE_NUM" +
		"\n" + "12345678,john,q,smith,555-555-5555" +
		"\n" + "55555555,,,,"

	const (
		wantErrorData = ERROR_HEADER + "\n" + "3," + string(errBadFirst)
		fileName      = "test.csv"
		inDir         = "test_in"
		outDir        = "test_out"
		errDir        = "test_err"
	)

	wantRecords := []Record{{
		ID:    12345678,
		First: "john", Middle: "q", Last: "smith",
		Phone: "555-555-5555",
	}}

	var d dirs
	var err error
	if d.in, err = ioutil.TempDir("", inDir); err != nil {
		t.Error(err)
	}

	defer os.RemoveAll(d.in)
	if d.out, err = ioutil.TempDir("", outDir); err != nil {
		t.Error(err)
	}

	defer os.RemoveAll(d.out)
	if d.err, err = ioutil.TempDir("", errDir); err != nil {
		t.Error(err)
	}

	ioutil.WriteFile(filepath.Join(d.in, fileName), []byte(testCSVData), 0666)

	os.Args = []string{"csv_to_json.exe", d.in, d.out, d.err}

	isTest = true
	main()
	gotJSONData, _ := ioutil.ReadFile(d.outPath(fileName))

	var gotRec []Record
	json.Unmarshal(gotJSONData, &gotRec)
	gotErrCSV, _ := ioutil.ReadFile(d.errPath(fileName))

	assert.Equal(t, wantRecords, gotRec)
	assert.Equal(t, wantErrorData, string(gotErrCSV))
}

func Test_getDirsFromArgs(t *testing.T) {
	want := dirs{"/foo", "/bar", "/baz"}
	args := []string{"csv_to_json", "/foo", "/bar", "/baz"}
	got, err := getDirsFromArgs(args)
	assert.Equal(t, want, got)
	assert.NoError(t, err)

	_, err = getDirsFromArgs(nil)
	assert.Error(t, err)
}

func Test_deleteCSV(t *testing.T) {
	const goodFile = "test_file_for_deletion.csv"
	ioutil.WriteFile(goodFile, nil, 0777)
	defer os.Remove(goodFile)
	assert.NoError(t, deleteCSV(goodFile))

	const badFile = "foo.baz"
	ioutil.WriteFile(badFile, nil, 0777)
	defer os.Remove(badFile)
	assert.Error(t, deleteCSV(badFile))

}
