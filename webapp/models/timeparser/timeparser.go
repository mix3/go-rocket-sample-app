package timeparser

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

type UnMatchError struct {
	error
}

type ParseError struct {
	error
}

type PassedError struct {
	error
}

var now = time.Now

var format = "2006-01-02 15:04:05 -0700"

var matchFunc = map[*regexp.Regexp]func(args ...int) (time.Time, error){
	regexp.MustCompile(`^(\d{1,2}):(\d{1,2})$`): func(args ...int) (time.Time, error) {
		str := fmt.Sprintf("%04d-%02d-%02d %02d:%02d:00 +0900", now().Year(), now().Month(), now().Day(), args[0], args[1])
		t, _ := time.Parse(format, str)
		if t.Before(now()) {
			return t.AddDate(0, 0, 1), nil
		}
		return t, nil
	},
	regexp.MustCompile(`^(\d{1,2})(?:/|-)(\d{1,2}) (\d{1,2}):(\d{1,2})$`): func(args ...int) (time.Time, error) {
		str := fmt.Sprintf("%04d-%02d-%02d %02d:%02d:00 +0900", now().Year(), args[0], args[1], args[2], args[3])
		t, err := time.Parse(format, str)
		if err != nil {
			return time.Unix(0, 0), ParseError{fmt.Errorf(`%v`, err)}
		}
		if t.Before(now()) {
			return time.Unix(0, 0), PassedError{fmt.Errorf(`%v is already passed`, t)}
		}
		return t, nil
	},
	regexp.MustCompile(`^(\d{4})(?:/|-)(\d{1,2})(?:/|-)(\d{1,2}) (\d{1,2}):(\d{1,2})$`): func(args ...int) (time.Time, error) {
		str := fmt.Sprintf("%04d-%02d-%02d %02d:%02d:00 +0900", args[0], args[1], args[2], args[3], args[4])
		t, err := time.Parse(format, str)
		if err != nil {
			return time.Unix(0, 0), ParseError{fmt.Errorf(`%v`, err)}
		}
		if t.Before(now()) {
			return time.Unix(0, 0), PassedError{fmt.Errorf(`%v is already passed`, t)}
		}
		return t, nil
	},
}

func Parse(timeStr string) (time.Time, error) {
	for r, f := range matchFunc {
		if r.MatchString(timeStr) {
			args := []int{}
			for _, v := range r.FindStringSubmatch(timeStr)[1:] {
				i, _ := strconv.Atoi(v)
				args = append(args, i)
			}
			return f(args...)
		}
	}
	return time.Unix(0, 0), UnMatchError{fmt.Errorf(`"%v" is invalid`, timeStr)}
}
