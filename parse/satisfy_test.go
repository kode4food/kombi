package parse_test

import (
	"fmt"
	"testing"

	"github.com/caravan/kombi/parse"
	"github.com/stretchr/testify/assert"
)

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
