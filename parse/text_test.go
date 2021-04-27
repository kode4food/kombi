package parse_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/caravan/kombi/parse"
	"github.com/stretchr/testify/assert"
)

func TestRegExp(t *testing.T) {
	as := assert.New(t)

	integer := parse.RegExp("[0-9]+").Map(func(r parse.Result) parse.Result {
		res, _ := strconv.ParseInt(r.(string), 0, 32)
		return int(res)
	})

	s, f := integer.Parse("1001")
	as.NotNil(s)
	as.Nil(f)
	as.Equal(1001, s.Result)

	s, f = integer.Parse("not")
	as.Nil(s)
	as.NotNil(f)
	as.EqualError(f.Error,
		fmt.Sprintf(parse.ErrWrappedExpectation,
			fmt.Sprintf(parse.ErrExpectedPattern, "[0-9]+"),
			"not",
		),
	)
}

func TestString(t *testing.T) {
	as := assert.New(t)

	strCmp := parse.String("Case Sensitive")
	s, f := strCmp.Parse("Case Sensitive")
	as.NotNil(s)
	as.Nil(f)
	as.Equal("Case Sensitive", s.Result)

	s, f = strCmp.Parse("CaSe SeNsItIve")
	as.Nil(s)
	as.NotNil(f)
	as.EqualError(f.Error,
		fmt.Sprintf(parse.ErrWrappedExpectation,
			fmt.Sprintf(parse.ErrExpectedString, "Case Sensitive"),
			"CaSe SeNsItIve",
		),
	)
}

func TestStrCaseCmp(t *testing.T) {
	as := assert.New(t)

	insCmp := parse.StrCaseCmp("Case Insensitive")
	s, f := insCmp.Parse("Case INSENSITIVE")
	as.NotNil(s)
	as.Nil(f)
	as.Equal("Case INSENSITIVE", s.Result)

	s, f = insCmp.Parse("Ca$e INSENSITIVE")
	as.Nil(s)
	as.NotNil(f)
	as.EqualError(f.Error,
		fmt.Sprintf(parse.ErrWrappedExpectation,
			fmt.Sprintf(parse.ErrExpectedString, "Case Insensitive"),
			"Ca$e INSENSITIVE",
		),
	)
}
