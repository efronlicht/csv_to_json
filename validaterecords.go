package main

import (
	"regexp"
	"strconv"
)

type validationError string

const (
	smallest8Digit = 10000000
	largest8Digit  = smallest8Digit*10 - 1

	//ID, First, Middle, Last, Phone
	FIELDS_PER_RECORD = 5

	errBadFirst            validationError = "first name must be between 1 and 15 characters"
	errBadMiddle           validationError = "middle name must be between 0 and 15 characters"
	errBadLast             validationError = "last name must be between 1 and 15 characters"
	errBadPhone            validationError = "phone must be in the form xxx-xxx-xxxx, where x are digits"
	errBadID               validationError = "UID must be an eight-digit positive integer"
	errWrongNumberOfFields validationError = "line must have exactly five fields: " +
		"INTERNAL_ID, FIRST_NAME, MIDDLE_NAME, LAST_NAME, PHONE_NUM"
)

type Record struct {
	ID     uint64 `json:"id"`
	First  string `json:"first"`
	Last   string `json:"last"`
	Middle string `json:"middle,omitempty"`
	Phone  string `json:"phone"`
}

func toRecord(fields []string) (rec Record, err error) {

	if len(fields) != FIELDS_PER_RECORD {
		return rec, errWrongNumberOfFields
	}

	id, _ := strconv.ParseUint(fields[0], 10, 64) //we validate later, so we can ignore this err
	var first, middle, last, phone = fields[1], fields[2], fields[3], fields[4]

	switch {
	case !validID(id):
		return rec, errBadID
	case !validFirst(first):
		return rec, errBadFirst
	case !validMiddle(middle):
		return rec, errBadMiddle
	case !validLast(last):
		return rec, errBadLast
	case !validPhone(phone):
		return rec, errBadPhone
	default:
		return Record{ID: id, First: first, Middle: middle, Last: last, Phone: phone}, nil
	}
}
func validFirst(name string) bool {
	return len(name) < 15 && len(name) > 0
}

var validLast = validFirst

func validMiddle(name string) bool {
	return len(name) < 15
}

var validPhoneRE = regexp.MustCompile("[0-9]{3}-[0-9]{3}-[0-9]{4}")

func validPhone(phone string) bool {
	return validPhoneRE.Match([]byte(phone))
}

func validID(id uint64) bool {
	return smallest8Digit <= id && id <= largest8Digit
}

func (err validationError) Error() string { return string(err) }
