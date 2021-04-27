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

	cmp := parse.String("Anything")
	s, f := cmp.Parse("Anything")
	as.NotNil(s)
	as.Nil(f)
	as.Equal("Anything", s.Result)

	s, f = cmp.Parse("aNyThiNg")
	as.Nil(s)
	as.NotNil(f)
	as.EqualError(f.Error,
		fmt.Sprintf(parse.ErrWrappedExpectation,
			fmt.Sprintf(parse.ErrExpectedString, "Anything"),
			"aNyThiNg",
		),
	)
}

func TestStrCaseCmp(t *testing.T) {
	as := assert.New(t)

	cmp := parse.StrCaseCmp("Anything")
	s, f := cmp.Parse("anyTHING")
	as.NotNil(s)
	as.Nil(f)
	as.Equal("anyTHING", s.Result)

	s, f = cmp.Parse("aNyThaNg")
	as.Nil(s)
	as.NotNil(f)
	as.EqualError(f.Error,
		fmt.Sprintf(parse.ErrWrappedExpectation,
			fmt.Sprintf(parse.ErrExpectedString, "Anything"),
			"aNyThaNg",
		),
	)
}
