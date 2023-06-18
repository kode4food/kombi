package parse_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/kode4food/kombi/parse"
)

func TestReturn(t *testing.T) {
	as := NewAssert(t)

	res := parse.Return("hello")
	s, f := res.Parse("this is a test")
	as.SuccessResult(s, f, "hello")
	as.Equal(parse.Input("this is a test"), s.Remaining)
}

func TestBind(t *testing.T) {
	as := NewAssert(t)

	integer := parse.RegExp("[0-9]+").Map(func(r parse.Result) parse.Result {
		res, _ := strconv.ParseInt(r.(string), 0, 32)
		return int(res)
	})

	add := integer.Bind(
		func(l parse.Result) parse.Parser {
			return parse.String("+").Bind(
				func(_ parse.Result) parse.Parser {
					return integer.Bind(
						func(r parse.Result) parse.Parser {
							return parse.Return(l.(int) + r.(int))
						},
					)
				},
			)
		},
	)

	s, f := add.Parse("2+8")
	as.SuccessResult(s, f, 10)
}

func TestCapture(t *testing.T) {
	as := NewAssert(t)

	var captured int
	integer := parse.RegExp("[0-9]+").Capture(func(r parse.Result) {
		res, _ := strconv.ParseInt(r.(string), 0, 32)
		captured = int(res)
	})

	s, f := integer.Parse("nope")
	as.FailureWrapped(s, f,
		fmt.Sprintf(parse.ErrExpectedPattern, "[0-9]+"), "nope",
	)
	as.Equal(0, captured)

	s, f = integer.Parse("42")
	as.SuccessResult(s, f, "42")
	as.Equal(42, captured)
}

func TestAnd(t *testing.T) {
	as := NewAssert(t)

	hello := parse.String("hello").EOF()
	s, f := hello.Parse("hello")
	as.Success(s, f)

	s, f = hello.Parse("hell no")
	as.FailureWrapped(s, f,
		fmt.Sprintf(parse.ErrExpectedString, "hello"),
		"hell no",
	)

	s, f = hello.Parse("hello you")
	as.FailureWrapped(s, f, parse.ErrExpectedEndOfFile, " you")
}

func TestOr(t *testing.T) {
	as := NewAssert(t)

	maybeHello := parse.EOF.Or(
		parse.String("hello").EOF(),
	)

	as.Success(maybeHello.Parse("hello"))
	as.Success(maybeHello.Parse(""))

	s, f := maybeHello.Parse("hello there")
	as.FailureWrapped(s, f, parse.ErrExpectedEndOfFile, " there")
}

func TestMap(t *testing.T) {
	as := NewAssert(t)

	intMapper := parse.RegExp("[0-9]+").Map(
		func(r parse.Result) parse.Result {
			if res, err := strconv.ParseInt(r.(string), 10, 32); err == nil {
				return int(res)
			}
			return 0
		},
	).Or(parse.Fail("couldn't parse int"))

	s, f := intMapper.Parse("42")
	as.SuccessResult(s, f, 42)

	s, f = intMapper.Parse("hello")
	as.FailureError(s, f, "couldn't parse int")
}
