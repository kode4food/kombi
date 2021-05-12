package parse_test

import (
	"errors"
	"testing"

	"github.com/caravan/kombi/parse"
)

func TestParserReturn(t *testing.T) {
	as := NewAssert(t)

	p := parse.String("hello").Return(true)
	s, f := p.Parse("hello")
	as.SuccessResult(s, f, true)
}

func TestParserFail(t *testing.T) {
	as := NewAssert(t)

	p := parse.String("hello").Fail("explode!")
	s, f := p.Parse("hello")
	as.FailureError(s, f, "explode!")
}

func TestParserSatisfy(t *testing.T) {
	as := NewAssert(t)

	p := parse.String("hello").Satisfy(func(i parse.Input) (int, error) {
		if i[0] == '!' {
			return 1, nil
		}
		return 0, errors.New("mismatch")
	})

	s, f := p.Parse("hello!")
	as.SuccessResult(s, f, parse.Input("!"))
}
