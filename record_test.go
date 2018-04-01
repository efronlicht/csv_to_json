package main

import (
	"strings"
	"testing"

	"github.com/eyecuelab/kit/copyslice"
	"github.com/stretchr/testify/assert"
)

func Test_toRecord(t *testing.T) {
	const (
		badID     = "0xff"
		badMiddle = "maskdmalsdkmklamsdlkamsd"
		badLast   = "lorem ipsum something something long string"
		badPhone  = "867309"

		goodFirst  = "john"
		goodMiddle = "q"
		goodLast   = "doe"
		goodPhone  = "555-555-5555"
		goodID     = "12345678"

		wantID = 12345678
	)

	var (
		goodFields          = []string{goodID, goodFirst, goodMiddle, goodLast, goodPhone}
		wantRecord          = Record{ID: wantID, First: goodFirst, Middle: goodMiddle, Last: goodLast, Phone: goodPhone}
		wrongNumberOfFields = []string{"foo", "bar"}
	)

	_, err := toRecord(wrongNumberOfFields)
	assert.Equal(t, err, errWrongNumberOfFields)

	record, _ := toRecord(goodFields)
	assert.Equal(t, record, wantRecord)

	wrongID := copyslice.String(goodFields)
	wrongID[0] = badID

	_, err = toRecord(wrongID)
	assert.Equal(t, err, errBadID)

	missingFirst := copyslice.String(goodFields)
	missingFirst[1] = ""
	_, err = toRecord(missingFirst)
	assert.Equal(t, err, errBadFirst)

	wrongMiddle := copyslice.String(goodFields)
	wrongMiddle[2] = strings.Repeat("foo", 12)
	_, err = toRecord(wrongMiddle)
	assert.Equal(t, errBadMiddle, err)

	wrongLast := copyslice.String(goodFields)
	wrongLast[3] = strings.Repeat("bar", 12)
	_, err = toRecord(wrongLast)
	assert.Equal(t, errBadLast, err)

	missingPhone := copyslice.String(goodFields)
	missingPhone[4] = ""
	_, err = toRecord(missingPhone)
	assert.Equal(t, err, errBadPhone)
}

func Test_validFirst_and_validLast(t *testing.T) {
	assert.True(t, validLast("Smith"))
	assert.False(t, validFirst(""))
	assert.False(t, validFirst(strings.Repeat("n", 16)))
}

func Test_validMiddle(t *testing.T) {
	assert.True(t, validMiddle(""))
	assert.True(t, validMiddle("mmmmmmmmmmmmmmm"))
	assert.False(t, validMiddle("the quick brown fox jumped over the lazy dog"))
}

func Test_validPhone(t *testing.T) {
	assert.True(t, validPhone("555-555-5555"))
	assert.False(t, validPhone("foobar"))
	assert.False(t, validPhone("1-555-281-9999"))
	assert.False(t, validPhone("9999999999"))
}

func Test_validID(t *testing.T) {
	assert.True(t, validID(smallest8Digit))
	assert.True(t, validID(largest8Digit))
	assert.False(t, validID(smallest8Digit-1))
	assert.False(t, validID(largest8Digit+1))

}
