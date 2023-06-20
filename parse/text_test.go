package parse_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/kode4food/kombi/parse"
)

func TestRegExp(t *testing.T) {
	as := NewAssert(t)

	integer := parse.RegExp("[0-9]+").Map(func(r any) any {
		res, _ := strconv.ParseInt(r.(string), 0, 32)
		return int(res)
	})

	s, f := integer.Parse("1001")
	as.SuccessResult(s, f, 1001)

	s, f = integer.Parse("not")
	as.FailureWrapped(s, f,
		fmt.Sprintf(parse.ErrExpectedPattern, "[0-9]+"), "not",
	)
}

func TestString(t *testing.T) {
	as := NewAssert(t)

	strCmp := parse.String("Case Sensitive")
	s, f := strCmp.Parse("Case Sensitive")
	as.SuccessResult(s, f, "Case Sensitive")

	s, f = strCmp.Parse("CaSe SeNsItIve")
	as.FailureWrapped(s, f,
		fmt.Sprintf(parse.ErrExpectedString, "Case Sensitive"),
		"CaSe SeNsItIve",
	)
}

func TestStrCaseCmp(t *testing.T) {
	as := NewAssert(t)

	insCmp := parse.StrCaseCmp("Case Insensitive")
	s, f := insCmp.Parse("Case INSENSITIVE")
	as.SuccessResult(s, f, "Case INSENSITIVE")

	s, f = insCmp.Parse("Ca$e INSENSITIVE")
	as.FailureWrapped(s, f,
		fmt.Sprintf(parse.ErrExpectedString, "Case Insensitive"),
		"Ca$e INSENSITIVE",
	)
}
