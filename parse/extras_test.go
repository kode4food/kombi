package parse_test

import (
	"fmt"
	"testing"

	"github.com/caravan/kombi/parse"
	"github.com/stretchr/testify/assert"
)

func TestAnyOf(t *testing.T) {
	as := assert.New(t)

	maybeGreet := parse.AnyOf(
		parse.String("hello").Then(parse.EOF),
		parse.String("howdy").Then(parse.EOF),
		parse.String("ciao").Then(parse.EOF),
		parse.EOF,
	)

	s, f := maybeGreet.Parse("hello")
	as.NotNil(s)
	as.Nil(f)

	s, f = maybeGreet.Parse("howdy")
	as.NotNil(s)
	as.Nil(f)

	s, f = maybeGreet.Parse("ciao")
	as.NotNil(s)
	as.Nil(f)

	s, f = maybeGreet.Parse("not")
	as.Nil(s)
	as.NotNil(f)
	as.EqualError(f.Error,
		fmt.Sprintf(parse.ErrWrappedExpectation,
			parse.ErrExpectedEndOfFile, "not",
		),
	)

	s, f = maybeGreet.Parse("way too long so will be truncated")
	as.Nil(s)
	as.NotNil(f)
	as.EqualError(f.Error,
		fmt.Sprintf(parse.ErrWrappedExpectation,
			parse.ErrExpectedEndOfFile, "way too long so ...",
		),
	)
}

func TestDefaulted(t *testing.T) {
	as := assert.New(t)

	optional := parse.String("hello").Optional()
	s, f := optional.Parse("hello")
	as.NotNil(s)
	as.Nil(f)
	as.Equal("hello", s.Result)
	as.Equal(parse.Input(""), s.Remaining)

	s, f = optional.Parse("doof")
	as.NotNil(s)
	as.Nil(f)
	as.Nil(s.Result)
	as.Equal(parse.Input("doof"), s.Remaining)

	defaulted := parse.String("hello").DefaultTo(func() parse.Result {
		return "nope"
	})
	s, f = defaulted.Parse("doof")
	as.NotNil(s)
	as.Nil(f)
	as.Equal("nope", s.Result)
	as.Equal(parse.Input("doof"), s.Remaining)
}

func TestIgnored(t *testing.T) {
	as := assert.New(t)

	ignored := parse.String("SKIP ").Ignore().Then(parse.String("THIS"))
	s, f := ignored.Parse("SKIP THIS")
	as.NotNil(s)
	as.Nil(f)
	as.Equal(1, len(s.Result.(parse.Results)))
	as.Equal("THIS", s.Result.(parse.Results)[0])

	s, f = ignored.Parse("NOT THIS")
	as.Nil(s)
	as.NotNil(f)
	as.EqualError(f.Error,
		fmt.Sprintf(parse.ErrWrappedExpectation,
			fmt.Sprintf(parse.ErrExpectedString, "SKIP "),
			"NOT THIS",
		),
	)
}
