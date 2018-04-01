package main

import (
	"encoding/csv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_errorsToCSV(t *testing.T) {
	const testErr validationError = "foo"

	errs := []csv.ParseError{
		{Line: 2, Err: testErr},
		{Line: 5, Err: testErr},
	}

	want := strings.Join([]string{ERROR_HEADER, "2,foo", "5,foo"}, "\n")
	assert.Equal(t, want, string(errorsToCSV(errs)))
}
