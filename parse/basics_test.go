package parse_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/caravan/kombi/parse"
	"github.com/stretchr/testify/assert"
)

func TestReturn(t *testing.T) {
	as := assert.New(t)

	res := parse.Return("hello")
	s, f := res.Parse("this is a test")
	as.NotNil(s)
	as.Nil(f)
	as.Equal("hello", s.Result)
	as.Equal(parse.Input("this is a test"), s.Remaining)
}

func TestBind(t *testing.T) {
	as := assert.New(t)

	b := func(r parse.Result) parse.Parser {
		if r.(string) == "hello" {
			return parse.String(" there!")
		}
		return parse.Fail("exploded")
	}

	bound := parse.AnyOf(
		parse.String("hello"),
		parse.String("explode"),
	).Bind(b)

	s, f := bound.Parse("hello there!")
	as.NotNil(s)
	as.Nil(f)
	as.Equal(2, len(s.Result.(parse.Results)))
	as.Equal("hello", s.Result.(parse.Results)[0])
	as.Equal(" there!", s.Result.(parse.Results)[1])

	s, f = bound.Parse("explode")
	as.Nil(s)
	as.NotNil(f)
	as.EqualError(f.Error, "exploded")
}

func TestAnd(t *testing.T) {
	as := assert.New(t)

	hello := parse.String("hello").Then(parse.EOF)
	s, f := hello.Parse("hello")
	as.NotNil(s)
	as.Nil(f)

	s, f = hello.Parse("hell no")
	as.Nil(s)
	as.NotNil(f)
	as.EqualError(f.Error,
		fmt.Sprintf(parse.ErrWrappedExpectation,
			fmt.Sprintf(parse.ErrExpectedString, "hello"),
			"hell no",
		),
	)

	s, f = hello.Parse("hello you")
	as.Nil(s)
	as.NotNil(f)
	as.EqualError(f.Error,
		fmt.Sprintf(parse.ErrWrappedExpectation,
			parse.ErrExpectedEndOfFile, " you",
		),
	)
}

func TestOr(t *testing.T) {
	as := assert.New(t)

	maybeHello := parse.EOF.Or(
		parse.String("hello").Then(parse.EOF),
	)

	s, f := maybeHello.Parse("hello")
	as.NotNil(s)
	as.Nil(f)

	s, f = maybeHello.Parse("")
	as.NotNil(s)
	as.Nil(f)

	s, f = maybeHello.Parse("hello there")
	as.Nil(s)
	as.NotNil(f)
	as.EqualError(f.Error,
		fmt.Sprintf(parse.ErrWrappedExpectation,
			parse.ErrExpectedEndOfFile, " there",
		),
	)
}

func TestMap(t *testing.T) {
	as := assert.New(t)

	intMapper := parse.RegExp("[0-9]+").Map(
		func(r parse.Result) parse.Result {
			if res, err := strconv.ParseInt(r.(string), 10, 32); err == nil {
				return int(res)
			}
			return 0
		},
	).Or(parse.Fail("couldn't parse int"))

	s, f := intMapper.Parse("42")
	as.NotNil(s)
	as.Nil(f)
	as.Equal(42, s.Result)

	s, f = intMapper.Parse("hello")
	as.Nil(s)
	as.NotNil(f)
	as.EqualError(f.Error, "couldn't parse int")
}
