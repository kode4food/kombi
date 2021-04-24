package parse_test

import (
	"bytes"
	"fmt"
	"strconv"
	"testing"

	"github.com/caravan/kombi/parse"
	"github.com/stretchr/testify/assert"
)

func TestAndThen(t *testing.T) {
	as := assert.New(t)

	hello := parse.String("hello").AndThen(parse.EndOfFile)
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

func TestOrElse(t *testing.T) {
	as := assert.New(t)

	maybeHello := parse.EndOfFile.OrElse(
		parse.String("hello").AndThen(parse.EndOfFile),
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

func TestAnyOf(t *testing.T) {
	as := assert.New(t)

	maybeGreet := parse.AnyOf(
		parse.String("hello").AndThen(parse.EndOfFile),
		parse.String("howdy").AndThen(parse.EndOfFile),
		parse.String("ciao").AndThen(parse.EndOfFile),
		parse.EndOfFile,
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

func TestMap(t *testing.T) {
	as := assert.New(t)

	intMapper := parse.RegExp("[0-9]+").Map(
		func(r parse.Result) parse.Result {
			if res, err := strconv.ParseInt(r.(string), 10, 32); err == nil {
				return int(res)
			}
			return 0
		},
	).OrElse(parse.Error("couldn't parse int"))

	s, f := intMapper.Parse("42")
	as.NotNil(s)
	as.Nil(f)
	as.Equal(42, s.Result)

	s, f = intMapper.Parse("hello")
	as.Nil(s)
	as.NotNil(f)
	as.EqualError(f.Error, "couldn't parse int")
}

func TestCombine(t *testing.T) {
	as := assert.New(t)

	helloThere := parse.
		String("hello ").AndThen(
		parse.String("there").AndThen(parse.String("!"))).
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

	goodbye := parse.AndThen(
		parse.String("good").AndThen(parse.String("")),
		parse.String("bye").AndThen(parse.String("")),
	).Combine(stringResults)

	s, f = goodbye.Parse("goodbye")
	as.NotNil(s)
	as.Nil(f)
	as.Equal("good->->bye->->", s.Result)
}

func stringResults(r ...parse.Result) parse.Result {
	var buf bytes.Buffer
	for _, e := range r {
		buf.WriteString(e.(string))
		buf.WriteString("->")
	}
	return buf.String()
}