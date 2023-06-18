package parse_test

import (
	"testing"

	"github.com/kode4food/kombi/parse"
)

func TestAny(t *testing.T) {
	as := NewAssert(t)

	maybeGreet := parse.Any(
		parse.String("hello").EOF(),
		parse.String("howdy").EOF(),
		parse.String("ciao").EOF(),
		parse.EOF,
	)

	as.Success(maybeGreet.Parse("hello"))
	as.Success(maybeGreet.Parse("howdy"))
	as.Success(maybeGreet.Parse("ciao"))

	s, f := maybeGreet.Parse("not")
	as.FailureWrapped(s, f, parse.ErrExpectedEndOfFile, "not")

	s, f = maybeGreet.Parse("way too long so will be truncated")
	as.FailureWrapped(s, f, parse.ErrExpectedEndOfFile, "way too long so ...")
}

func TestDefaulted(t *testing.T) {
	as := NewAssert(t)

	optional := parse.String("hello").Optional()
	s, f := optional.Parse("hello")
	as.SuccessResult(s, f, "hello")
	as.Equal(parse.Input(""), s.Remaining)

	s, f = optional.Parse("doof")
	as.SuccessResult(s, f, nil)
	as.Equal(parse.Input("doof"), s.Remaining)

	defaulted := parse.String("hello").DefaultTo("nope")
	s, f = defaulted.Parse("doof")
	as.SuccessResult(s, f, "nope")
	as.Equal(parse.Input("doof"), s.Remaining)
}
