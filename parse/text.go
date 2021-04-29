package parse

import (
	"regexp"
	"strings"
)

type mapString func(string) string

// Error messages
const (
	ErrExpectedPattern = "expected pattern: %s"
	ErrExpectedString  = "expected string %s"
)

// RegExp returns a Parser that is used to Satisfy an IsRegExp Predicate
func RegExp(s string) Parser {
	return Satisfy(IsRegExp(s)).Map(toString)
}

// IsRegExp returns a Predicate that can be used to Satisfy regular
// expression patterns in the Input
func IsRegExp(s string) Predicate {
	pattern := regexp.MustCompile("^(" + s + ")")
	return func(i Input) (int, error) {
		if sm := pattern.FindStringSubmatch(string(i)); sm != nil {
			matched := sm[0]
			return len(matched), nil
		}
		return 0, i.errExpected(ErrExpectedPattern, s)
	}
}

// String returns a Parser that is used to Satisfy an IsString Predicate
func String(s string) Parser {
	return Satisfy(IsString(s)).Map(toString)
}

// IsString returns a Predicate that can be used to Satisfy case-sensitive
// string patterns in the Input
func IsString(s string) Predicate {
	return stringPredicate(s, stringIdentity)
}

// StrCaseCmp returns a Parser that is used to Satisfy an IsStrCaseCmp
// Predicate
func StrCaseCmp(s string) Parser {
	return Satisfy(IsStrCaseCmp(s)).Map(toString)
}

// IsStrCaseCmp returns a Predicate that can be used to Satisfy
// case-insensitive string patterns in the Input
func IsStrCaseCmp(s string) Predicate {
	return stringPredicate(s, strings.ToUpper)
}

func stringIdentity(s string) string {
	return s
}

func stringPredicate(s string, normalize mapString) Predicate {
	n := normalize(s)
	size := len(n)
	return func(i Input) (int, error) {
		if len(i) >= size {
			cmp := string(i[0:size])
			if n == normalize(cmp) {
				return len(cmp), nil
			}
		}
		return 0, i.errExpected(ErrExpectedString, s)
	}
}

func toString(r Result) Result {
	return string(r.(Input))
}
