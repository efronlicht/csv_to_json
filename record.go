package main

import (
	"regexp"
	"strconv"
)

type validationError string

const (
	errBadFirst            validationError = "missing or invalid first name: must be between 1 and 15 characters"
	errBadMiddle           validationError = "invalid middle name: must be between 0 and 15 characters"
	errBadLast             validationError = "missing or invalid last name: must be between 1 and 15 characters"
	errBadPhone            validationError = "missing or invalid phone #: must be in the form xxx-xxx-xxxx, where x are digits"
	errBadID               validationError = "missing or invalid id: must be an eight-digit positive integer"
	errWrongNumberOfFields validationError = "wrong number of fields: line must have exactly five fields: " +
		"INTERNAL_ID, FIRST_NAME, MIDDLE_NAME, LAST_NAME, PHONE_NUM"

	smallest8Digit = 10000000
	largest8Digit  = smallest8Digit*10 - 1

	//ID, First, Middle, Last, Phone
	FIELDS_PER_RECORD = 5
)

type Record struct {
	ID     uint64 `json:"id"`
	First  string `json:"first"`
	Middle string `json:"middle,omitempty"`
	Last   string `json:"last"`
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
		rec = Record{
			ID:    id,
			First: first, Middle: middle, Last: last,
			Phone: phone,
		}
		return rec, nil
	}
}
func validFirst(name string) bool {
	return len(name) > 0 && len(name) <= 15
}

var validLast = validFirst

func validMiddle(name string) bool {
	return len(name) <= 15
}

var validPhoneRE = regexp.MustCompile("^[0-9]{3}-[0-9]{3}-[0-9]{4}$")

func validPhone(phone string) bool {
	return validPhoneRE.Match([]byte(phone))
}

func validID(id uint64) bool {
	return smallest8Digit <= id && id <= largest8Digit
}

func (err validationError) Error() string { return string(err) }
