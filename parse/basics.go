package parse

type (
	// Binder returns a Parser based on the provided result
	Binder func(any) Parser

	// Accept receives a result from a Capture Parser
	Accept func(any)

	// Mapper maps one result value to another
	Mapper func(any) any

	// Predicate checks the beginning of its provided Input for a match
	Predicate func(Input) (int, error)

	eof struct{}
)

// Error messages
const (
	ErrExpectedEndOfFile = "expected end of file"
)

// EndOfFile represents the matched EOF result
var EndOfFile = &eof{}

// Return returns a new Parser. This Parser consumes none of the Input, but
// instead returns a Success containing the provided result
func Return(r any) Parser {
	return func(i Input) (*Success, *Failure) {
		return i.succeedWith(r)
	}
}

// Bind returns a new Parser, the result of which is based on the result of the
// provided Parser being Combined with the result of the Parser returned by the
// provided Binder
func Bind(p Parser, b Binder) Parser {
	return func(i Input) (*Success, *Failure) {
		s, f := p(i)
		if f == nil {
			return b(s.Result)(s.Remaining)
		}
		return nil, f
	}
}

// Capture returns a new Parser, the result of which is taken from the supplied
// Parser and provided to the Accept function
func Capture(p Parser, accept Accept) Parser {
	return Map(p, func(r any) any {
		accept(r)
		return r
	})
}

// Map returns a new Parser, the result of which is a value generated by the
// provided Mapper
func Map(p Parser, fn Mapper) Parser {
	return Bind(p, func(r any) Parser {
		return Return(fn(r))
	})
}

// Fail returns a Parser node that generates the specified error
func Fail(msg string, args ...any) Parser {
	return func(i Input) (*Success, *Failure) {
		return i.failMessage(msg, args...)
	}
}

// Satisfy returns a new Parser. This Parser consumes enough of the Input to
// satisfy the provided Predicate and returns Success on a match
func Satisfy(p Predicate) Parser {
	return func(i Input) (*Success, *Failure) {
		m, err := p(i)
		if err == nil {
			return i.succeedMatch(m)
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

// Then returns a new Parser based on the result of the left Parser being
// Combined with the results of the right Parser
func Then(l Parser, r Parser) Parser {
	return Bind(l, func(_ any) Parser {
		return r
	})
}

// Or returns a new Parser based on either the successful result of the left
// Parser or the result of the right Parser
func Or(l Parser, r Parser) Parser {
	return func(i Input) (*Success, *Failure) {
		if s, f := l(i); f == nil {
			return s, nil
		}
		return r(i)
	}
}
