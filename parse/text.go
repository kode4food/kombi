package parse

import (
	"regexp"
	"strings"
)

// Error messages
const (
	ErrExpectedPattern = "expected pattern: %s"
	ErrExpectedString  = "expected string %s"
)

// RegExp returns a Parser node that performs regular expression
// matching at the beginning of its Input
func RegExp(s string) Parser {
	pattern := regexp.MustCompile("^(" + s + ")")
	return Satisfy(func(i Input) (int, error) {
		if sm := pattern.FindStringSubmatch(string(i)); sm != nil {
			matched := sm[0]
			return len(matched), nil
		}
		return 0, i.errExpected(ErrExpectedPattern, s)
	}).Map(toString)
}

// String returns a Parser that matches the string provided to it. The
// resulting Parser performs case-sensitive matching
func String(s string) Parser {
	return stringParser(s, stringIdentify)
}

// StrCaseCmp returns a Parser that matches the string provided to it. The
// resulting Parser performs case-insensitive matching
func StrCaseCmp(s string) Parser {
	return stringParser(s, strings.ToUpper)
}

func stringIdentify(s string) string {
	return s
}

func stringParser(s string, normalize func(string) string) Parser {
	n := normalize(s)
	size := len(n)
	return Satisfy(func(i Input) (int, error) {
		if len(i) >= size {
			cmp := string(i[0:size])
			if n == normalize(cmp) {
				return len(cmp), nil
			}
		}
		return 0, i.errExpected(ErrExpectedString, s)
	}).Map(toString)
}

func toString(r Result) Result {
	return string(r.(Input))
}
