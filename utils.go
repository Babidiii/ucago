package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var longMonthNames = []string{
	"janvier",
	"février",
	"mars",
	"avril",
	"mai",
	"juin",
	"juillet",
	"août",
	"septembre",
	"octobre",
	"novembre",
	"décembre",
}

var errBad = errors.New("bad value for field")

func lookup(tab []string, val string) (string, error) {
	for i, v := range tab {
		if len(val) >= len(v) && val[0:len(v)] == v {
			return fmt.Sprintf("%.2d", i+1), nil
		}
	}
	return "", errBad
}

func ParseMonth(val string) (string, error) {
	return lookup(longMonthNames, val)
}

func GetDateFormat(val string) string {
	vals := strings.Split(val, " ")
	month, _ := ParseMonth(vals[2])
	i, _ := strconv.Atoi(vals[1])

	return fmt.Sprintf("%.2d/%s", i, month)
}

func lpad(s string, pad string, plength int) string {
	for i := len(s); i < plength; i++ {
		s = pad + s
	}
	return s
}
