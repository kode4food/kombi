package parse_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/caravan/kombi/parse"
	"github.com/stretchr/testify/assert"
)

func TestCombine(t *testing.T) {
	as := assert.New(t)

	helloThere := parse.
		String("hello ").Then(
		parse.String("there").Then(parse.String("!"))).
		Combine(stringResults)

	s, f := helloThere.Parse("hello there!")
	as.NotNil(s)
	as.Nil(f)
	as.Equal("hello ->there->!->", s.Result)

	hello := parse.String("hello").
		Combine(func(r ...parse.Result) parse.Result {
			return fmt.Sprintf("{%s}", r[0].(string))
		})
	s, f = hello.Parse("hello")
	as.NotNil(s)
	as.Nil(f)
	as.Equal("{hello}", s.Result)

	goodbye := parse.Then(
		parse.String("good").Then(parse.String("")),
		parse.String("bye").Then(parse.String("")),
	).Combine(stringResults)

	s, f = goodbye.Parse("goodbye")
	as.NotNil(s)
	as.Nil(f)
	as.Equal("good->->bye->->", s.Result)
}

func TestOneOrMore(t *testing.T) {
	as := assert.New(t)

	many := parse.String("hello").OneOrMore()
	s, f := many.Parse("hellohellohello")
	as.NotNil(s)
	as.Nil(f)
	as.Equal(3, len(s.Result.(parse.Results)))
	as.Equal("hello", s.Result.(parse.Results)[2])

	s, f = many.Parse("blah")
	as.Nil(s)
	as.NotNil(f)
	as.EqualError(f.Error,
		fmt.Sprintf(parse.ErrWrappedExpectation,
			fmt.Sprintf(parse.ErrExpectedString, "hello"),
			"blah",
		),
	)
}

func TestZeroOrMore(t *testing.T) {
	as := assert.New(t)

	many := parse.String("hello").ZeroOrMore()
	s, f := many.Parse("hellohellohello")
	as.NotNil(s)
	as.Nil(f)
	as.Equal(3, len(s.Result.(parse.Results)))
	as.Equal("hello", s.Result.(parse.Results)[2])

	s, f = many.Parse("blah")
	as.NotNil(s)
	as.Nil(f)
	as.Equal(0, len(s.Result.(parse.Results)))
	as.Equal(parse.Input("blah"), s.Remaining)
}

func stringResults(r ...parse.Result) parse.Result {
	var buf bytes.Buffer
	for _, e := range r {
		buf.WriteString(e.(string))
		buf.WriteString("->")
	}
	return buf.String()
}
