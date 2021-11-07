package parse

type (
	// Parser is the signature for a parsing node
	Parser func(Input) (*Success, *Failure)

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

	// Result represents a Parser's Success result
	Result interface{}
)

// Parse uses the current Parser to match the provided string
func (p Parser) Parse(s string) (*Success, *Failure) {
	return p(Input(s))
}

// Return returns a new Parser. This Parser consumes none of the Input, but
// instead returns a Success containing the provided Result
func (p Parser) Return(r Result) Parser {
	return p.Then(Return(r))
}

// Bind returns a new Parser, the Result of which is based on the Result of
// this Parser being Combined with the Result of the Parser returned by the
// provided Binder
func (p Parser) Bind(b Binder) Parser {
	return Bind(p, b)
}

// Capture returns a new Parser, the Result of which is taken from this Parser
// and provided to the Accept function
func (p Parser) Capture(a Accept) Parser {
	return Capture(p, a)
}

// Map returns a new Parser, the Result of which is a value generated by the
// provided Mapper
func (p Parser) Map(fn Mapper) Parser {
	return Map(p, fn)
}

// Fail returns a Parser node that generates the specified error
func (p Parser) Fail(msg string, args ...interface{}) Parser {
	return p.Then(Fail(msg, args...))
}

// Satisfy returns a new Parser. This Parser consumes enough of the Input to
// satisfy the provided Predicate and returns Success on a match
func (p Parser) Satisfy(pred Predicate) Parser {
	return p.Then(Satisfy(pred))
}

// EOF matches the end of the Input
func (p Parser) EOF() Parser {
	return p.Then(EOF)
}

// Then returns a new Parser based on the Result of this Parser being Combined
// with the results of the other Parser
func (p Parser) Then(other Parser) Parser {
	return Then(p, other)
}

// Or returns a new Parser based on either the successful Result of this Parser
// or the Result of the other Parser
func (p Parser) Or(other Parser) Parser {
	return Or(p, other)
}

// Optional returns a new Parser that will DefaultTo nil if the match is not
// successful
func (p Parser) Optional() Parser {
	return Optional(p)
}

// DefaultTo returns a new Parser that will return the provided Result if this
// Parser match is not successful
func (p Parser) DefaultTo(r Result) Parser {
	return DefaultTo(p, r)
}

// Concat returns a new Parser, the Result of which is generated by
// concatenating the Results of the provided Parsers
func (p Parser) Concat(other Parser) Parser {
	return Concat(p, other)
}

// Combine returns a new Parser, the Result of which is a value generated by
// passing any Combined results to the provided Combiner
func (p Parser) Combine(fn Combiner) Parser {
	return Combine(p, fn)
}

// OneOrMore returns a new Parser, the Result of which is the Combined set of
// values matched by the provided Parser being performed one or more times
func (p Parser) OneOrMore() Parser {
	return OneOrMore(p)
}

// ZeroOrMore returns a new Parser, the Result of which is the Combined set of
// values matched by the provided Parser being performed zero or more times
func (p Parser) ZeroOrMore() Parser {
	return ZeroOrMore(p)
}

// Delimited returns a new Parser, the Result of which is the Combined set of
// values matched by the provided Parser and delimited by the provided
// Delimiter, performed one or more times
func (p Parser) Delimited(d Delimiter) Parser {
	return Delimited(p, d)
}
