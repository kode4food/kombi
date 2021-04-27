package parse

type (
	// Binder returns a Parser based on the provided Result
	Binder func(Result) Parser

	// Mapper maps one Result value to another
	Mapper func(Result) Result

	// Predicate checks the beginning of its provided Input for a match
	Predicate func(Input) (int, error)

	// Result represents a Parser's Success result
	Result interface{}

	// Results represents multiple Results that have been combined. This
	// is usually the result of the Bind or Then combinators
	Results []Result

	eof  struct{}
	skip struct{}
)

// Error messages
const (
	ErrExpectedEndOfFile = "expected end of file"
)

var (
	// EndOfFile represents the matched EOF Result
	EndOfFile = &eof{}

	// Skip represents a Parser Result that should be ignored
	Skip = &skip{}
)

// Return returns a new Parser. This Parser consumes none of the Input,
// but instead returns a Success containing the provided Result
func Return(r Result) Parser {
	return func(i Input) (*Success, *Failure) {
		return i.succeedWith(r)
	}
}

// Bind returns a new Parser, the Result of which is based on the Result
// of the provided Parser being Combined with the Result of the Parser
// returned by the provided Binder
func Bind(p Parser, b Binder) Parser {
	return func(i Input) (*Success, *Failure) {
		if ls, f := p(i); f != nil {
			return nil, f
		} else if rs, f := b(ls.Result)(ls.Remaining); f != nil {
			return i.failThrough(f)
		} else {
			res := combineResults(ls.Result, rs.Result)
			return rs.Remaining.succeedWith(res)
		}
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

// Fail returns a Parser node that generates the specified error
func Fail(msg string, args ...interface{}) Parser {
	return func(i Input) (*Success, *Failure) {
		return i.failMessage(msg, args...)
	}
}

// Satisfy returns a new Parser. This Parser consumes enough of the Input
// to satisfy the provided Predicate and returns Success on a match
func Satisfy(p Predicate) Parser {
	return func(i Input) (*Success, *Failure) {
		m, err := p(i)
		if err == nil {
			return i[m:].succeedWith(i[0:m])
		}
		return i.failWith(err)
	}
}

// EOF is a Parser that matches the end of the Input
var EOF = Parser(func(i Input) (*Success, *Failure) {
	if len(i) == 0 {
		return i.succeedWith(EndOfFile)
	}
	return i.failExpected(ErrExpectedEndOfFile)
})

// Then returns a new Parser based on the Result of the left Parser being
// Combined with the results of the right Parser
func Then(l Parser, r Parser) Parser {
	return Bind(l, func(_ Result) Parser {
		return r
	})
}

// Or returns a new Parser based on either the successful Result of the
// left Parser or the Result of the right Parser
func Or(l Parser, r Parser) Parser {
	return func(i Input) (*Success, *Failure) {
		if s, f := l(i); f == nil {
			return s, nil
		}
		return r(i)
	}
}

func combineResults(l, r Result) Result {
	var res Results
	res = appendResults(res, l)
	res = appendResults(res, r)
	return res
}

func appendResults(res Results, r Result) Results {
	c, ok := r.(Results)
	if !ok {
		return appendResults(res, Results{r})
	}
	for _, e := range c {
		if _, ok := e.(*skip); !ok {
			res = append(res, e)
		}
	}
	return res
}
