package parse

type (
	// Combiner takes multiple result values and combines them into one
	Combiner func(...any) any

	// Delimiter is a Parser, but will be ignored when combined
	Delimiter = Parser

	// Results represents multiple Results that have been combined. This
	// is usually the result of the Bind or Then combinators
	Results []any
)

// Concat returns a new Parser, the result of which is generated by
// concatenating the Results of the provided Parsers
func Concat(l Parser, r Parser) Parser {
	return l.Bind(func(lr any) Parser {
		return r.Bind(func(rr any) Parser {
			return Return(concatResults(lr, rr))
		})
	})
}

// Combine returns a new Parser, the result of which is a value generated by
// passing any Combined results to the provided Combiner
func Combine(p Parser, fn Combiner) Parser {
	return p.Map(func(r any) any {
		if res, ok := r.(Results); ok {
			return fn(res...)
		}
		return fn(r)
	})
}

// OneOrMore returns a new Parser, the result of which is the Combined set of
// values matched by the provided Parser being performed one or more times
func OneOrMore(p Parser) Parser {
	return Concat(p, ZeroOrMore(p))
}

// ZeroOrMore returns a new Parser, the result of which is the Combined set of
// values matched by the provided Parser being performed zero or more times
func ZeroOrMore(p Parser) Parser {
	return Or(
		p.Bind(func(r any) Parser {
			return Concat(Return(r), ZeroOrMore(p))
		}),
		Return(Results{}),
	)
}

// Delimited returns a new Parser, the result of which is the Combined set of
// values matched by the provided Parser and delimited by the provided
// Delimiter, performed one or more times
func Delimited(p Parser, d Delimiter) Parser {
	return Concat(p, ZeroOrMore(d.Then(p)))
}

func concatResults(l, r any) Results {
	var res Results
	res = appendResults(res, l)
	res = appendResults(res, r)
	return res
}

func appendResults(res Results, r any) Results {
	if c, ok := r.(Results); ok {
		return append(res, c...)
	}
	return append(res, r)
}
