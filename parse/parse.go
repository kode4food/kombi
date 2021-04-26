package parse

import (
	"fmt"
	"regexp"
	"strings"
)

type (
	// Parser is the signature for a parsing node
	Parser func(Input) (*Success, *Failure)

	// Mapper maps one Result value to another
	Mapper func(Result) Result

	// Combiner takes multiple Result values and combines them into one
	Combiner func(...Result) Result

	// Emitter returns a Result
	Emitter func() Result

	// Input represents a Parser's input
	Input string

	// Result represents a Parser's Success result
	Result interface{}

	// Combined represents multiple Results that have been combined.
	// This is usually the result of the AndThen combinator
	Combined []Result

	// Success is the structure returned if the Parser is able to
	// successfully match its Input. Remaining is what remains unparsed
	// from the original Input value
	Success struct {
		Result
		Remaining Input
	}

	// Failure is the structure returned if the Parser is not able to
	// successfully match its Input
	Failure struct {
		Error error
		Input
	}

	arg  = interface{}
	eof  struct{}
	skip struct{}
)

// Error messages
const (
	ErrExpectedPattern   = "expected pattern: %s"
	ErrExpectedString    = "expected string %s"
	ErrExpectedEndOfFile = "expected end of file"

	ErrWrappedExpectation = "%s, got %s"
)

const (
	maxExpectedGot = 16
)

var (
	// EOF represents the matched EndOfFile Result
	EOF = &eof{}

	// Skip represents a Parser Result that should be ignored
	Skip = &skip{}
)

// Parse uses the current Parser to match the provided string
func (p Parser) Parse(s string) (*Success, *Failure) {
	return p(Input(s))
}

// AndThen returns a new Parser based on the Result of this Parser
// being Combined with the results of the other Parser
func (p Parser) AndThen(other Parser) Parser {
	return AndThen(p, other)
}

// OrElse returns a new Parser based on either the successful Result of
// this Parser or the Result of the other Parser
func (p Parser) OrElse(other Parser) Parser {
	return OrElse(p, other)
}

// Optional returns a new Parser that will DefaultTo nil if the match
// is not successful
func (p Parser) Optional() Parser {
	return Optional(p)
}

// DefaultTo returns a new Parser that will return the Result provided
// by the Emitter if the match is not successful
func (p Parser) DefaultTo(e Emitter) Parser {
	return DefaultTo(p, e)
}

// Ignore returns a new Parser, the result of which is ignored if
// matching is successful
func (p Parser) Ignore() Parser {
	return Ignore(p)
}

// Map returns a new Parser, the Result of which is a value generated
// by the provided Mapper
func (p Parser) Map(fn Mapper) Parser {
	return Map(p, fn)
}

// Combine returns a new Parser, the Result of which is a value
// generated by passing any Combined results to the provided Combiner
func (p Parser) Combine(fn Combiner) Parser {
	return Combine(p, fn)
}

// OneOrMore returns a new Parser, the Result of which is the Combined
// set of values matched by the provided Parser being performed one or
// more times
func (p Parser) OneOrMore() Parser {
	return OneOrMore(p)
}

// ZeroOrMore returns a new Parser, the Result of which is the Combined
// set of values matched by the provided Parser being performed zero or
// more times
func (p Parser) ZeroOrMore() Parser {
	return ZeroOrMore(p)
}

// AnyOf returns a Parser, the result of which is generated by
// attempting the provided Parsers in succession. The first Parser that
// returns a Success ends the processing, and its Success instance is
// returned
func AnyOf(parsers ...Parser) Parser {
	return func(i Input) (*Success, *Failure) {
		var s *Success
		var f *Failure
		for _, p := range parsers {
			if s, f = p(i); s != nil {
				return s, nil
			}
		}
		return nil, f
	}
}

// String returns a Parser that matches the string provided to it. The
// resulting Parser performs case-sensitive matching
func String(s string) Parser {
	l := len(s)
	return func(i Input) (*Success, *Failure) {
		if len(i) >= l {
			cmp := string(i[0:l])
			if s == cmp {
				return i[l:].succeedWith(cmp)
			}
		}
		return i.failWithExpected(ErrExpectedString, s)
	}
}

// StrCaseCmp returns a Parser that matches the string provided to it. The
// resulting Parser performs case-insensitive matching
func StrCaseCmp(s string) Parser {
	su := strings.ToUpper(s)
	l := len(su)
	return func(i Input) (*Success, *Failure) {
		if len(i) >= l {
			cmp := string(i[0:l])
			if su == strings.ToUpper(cmp) {
				return i[l:].succeedWith(cmp)
			}
		}
		return i.failWithExpected(ErrExpectedString, s)
	}
}

// Error returns a Parser node that generates the specified Error
func Error(msg string, args ...interface{}) Parser {
	return func(i Input) (*Success, *Failure) {
		return i.failWith(msg, args...)
	}
}

// RegExp returns a Parser node that performs regular expression
// matching at the beginning of its Input
func RegExp(s string) Parser {
	p := regexp.MustCompile("^(" + s + ")")
	return func(i Input) (*Success, *Failure) {
		src := string(i)
		if sm := p.FindStringSubmatch(src); sm != nil {
			matched := sm[0]
			return Input(src[len(matched):]).succeedWith(matched)
		}
		return i.failWithExpected(ErrExpectedPattern, s)
	}
}

func combineResults(l, r Result) Result {
	var res Combined
	res = appendResults(res, l)
	res = appendResults(res, r)
	return res
}

func appendResults(res Combined, r Result) Combined {
	c, ok := r.(Combined)
	if !ok {
		return appendResults(res, Combined{r})
	}
	for _, e := range c {
		if _, ok := e.(*skip); !ok {
			res = append(res, e)
		}
	}
	return res
}

// Combine returns a new Parser, the Result of which is a value
// generated by passing any Combined results to the provided Combiner
func Combine(p Parser, fn Combiner) Parser {
	return p.Map(func(r Result) Result {
		res := appendResults(Combined{}, r)
		return fn(res...)
	})
}

// AndThen returns a new Parser based on the Result of the left Parser
// being Combined with the results of the right Parser
func AndThen(l Parser, r Parser) Parser {
	return func(i Input) (*Success, *Failure) {
		if ls, f := l(i); f != nil {
			return nil, f
		} else if rs, f := r(ls.Remaining); f != nil {
			return i.failThrough(f)
		} else {
			res := combineResults(ls.Result, rs.Result)
			return rs.Remaining.succeedWith(res)
		}
	}
}

// OrElse returns a new Parser based on either the successful Result of
// the left Parser or the Result of the right Parser
func OrElse(l Parser, r Parser) Parser {
	return func(i Input) (*Success, *Failure) {
		if s, f := l(i); f == nil {
			return s, nil
		}
		return r(i)
	}
}

// Map returns a new Parser, the Result of which is a value generated
// by the provided Mapper
func Map(p Parser, fn Mapper) Parser {
	return func(i Input) (*Success, *Failure) {
		s, f := p(i)
		if f == nil {
			return s.Remaining.succeedWith(fn(s.Result))
		}
		return nil, f
	}
}

// Optional returns a new Parser that will DefaultTo nil if the match
// is not successful
func Optional(p Parser) Parser {
	return DefaultTo(p, func() Result {
		return nil
	})
}

// DefaultTo returns a new Parser that will return the Result provided
// by the Emitter if the match is not successful
func DefaultTo(p Parser, missing Emitter) Parser {
	return func(i Input) (*Success, *Failure) {
		if s, f := p(i); f == nil {
			return s, nil
		}
		return i.succeedWith(missing())
	}
}

// Ignore returns a new Parser, the result of which is ignored if
// matching is successful
func Ignore(p Parser) Parser {
	return func(i Input) (*Success, *Failure) {
		s, f := p(i)
		if f == nil {
			return s.Remaining.succeedWith(Skip)
		}
		return nil, f
	}
}

// OneOrMore returns a new Parser, the Result of which is the Combined
// set of values matched by the provided Parser being performed one or
// more times
func OneOrMore(p Parser) Parser {
	return nOrMore(1, p)
}

// ZeroOrMore returns a new Parser, the Result of which is the Combined
// set of values matched by the provided Parser being performed zero or
// more times
func ZeroOrMore(p Parser) Parser {
	return nOrMore(0, p)
}

func nOrMore(min int, p Parser) Parser {
	return func(i Input) (*Success, *Failure) {
		var res Combined
		next := i
		for count := 0; ; count++ {
			s, f := p(next)
			if f == nil {
				res = append(res, s.Result)
				next = s.Remaining
				continue
			}
			if count >= min {
				return next.succeedWith(res)
			}
			return nil, f
		}
	}
}

// EndOfFile is a Parser that matches the end of the Input
var EndOfFile = Parser(func(i Input) (*Success, *Failure) {
	if len(i) == 0 {
		return i.succeedWith(EOF)
	}
	return i.failWithExpected(ErrExpectedEndOfFile)
})

func (i Input) succeedWith(r Result) (*Success, *Failure) {
	return &Success{
		Result:    r,
		Remaining: i,
	}, nil
}

func (i Input) failWithExpected(msg string, args ...arg) (*Success, *Failure) {
	got := i
	if len(got) > maxExpectedGot {
		got = got[0:maxExpectedGot] + "..."
	}
	errMsg := fmt.Sprintf(msg, args...)
	err := fmt.Errorf(ErrWrappedExpectation, errMsg, got)
	return nil, &Failure{
		Error: err,
		Input: i,
	}
}

func (i Input) failWith(msg string, args ...arg) (*Success, *Failure) {
	return nil, &Failure{
		Error: fmt.Errorf(msg, args...),
		Input: i,
	}
}

func (i Input) failThrough(f *Failure) (*Success, *Failure) {
	return nil, &Failure{
		Error: f.Error,
		Input: i,
	}
}
