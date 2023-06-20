package parse_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/kode4food/kombi/parse"
)

func TestCombine(t *testing.T) {
	as := NewAssert(t)

	helloThere := parse.
		String("hello ").Concat(
		parse.String("there").Concat(parse.String("!"))).
		Combine(stringResults)

	s, f := helloThere.Parse("hello there!")
	as.SuccessResult(s, f, "hello ->there->!->")

	hello := parse.String("hello").
		Combine(func(r ...any) any {
			return fmt.Sprintf("{%s}", r[0].(string))
		})
	s, f = hello.Parse("hello")
	as.SuccessResult(s, f, "{hello}")

	goodbye := parse.Concat(
		parse.String("good").Concat(parse.String("")),
		parse.String("bye").Concat(parse.String("")),
	).Combine(stringResults)

	s, f = goodbye.Parse("goodbye")
	as.SuccessResult(s, f, "good->->bye->->")
}

func TestOneOrMore(t *testing.T) {
	as := NewAssert(t)

	many := parse.String("hello").OneOrMore()
	s, f := many.Parse("hellohellohello")
	as.SuccessResults(s, f, "hello", "hello", "hello")

	s, f = many.Parse("blah")
	as.FailureWrapped(s, f,
		fmt.Sprintf(parse.ErrExpectedString, "hello"),
		"blah",
	)
}

func TestZeroOrMore(t *testing.T) {
	as := NewAssert(t)

	many := parse.String("hello").ZeroOrMore()
	s, f := many.Parse("hellohellohello")
	as.SuccessResults(s, f, "hello", "hello", "hello")

	s, f = many.Parse("blah")
	as.SuccessResults(s, f)
	as.Equal(parse.Input("blah"), s.Remaining)
}

func TestDelimited(t *testing.T) {
	as := NewAssert(t)

	nums := parse.RegExp("[0-9]+").Delimited(parse.String(","))
	s, f := nums.Parse("1,2,42")
	as.SuccessResults(s, f, "1", "2", "42")
}

func stringResults(r ...any) any {
	var buf bytes.Buffer
	for _, e := range r {
		buf.WriteString(e.(string))
		buf.WriteString("->")
	}
	return buf.String()
}
