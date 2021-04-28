package parse_test

import (
	"errors"
	"testing"

	"github.com/caravan/kombi/parse"
	"github.com/stretchr/testify/assert"
)

func TestParserReturn(t *testing.T) {
	as := assert.New(t)

	p := parse.String("hello").Return(true)
	s, f := p.Parse("hello")
	as.NotNil(s)
	as.Nil(f)
	as.True(s.Result.(bool))
}

func TestParserFail(t *testing.T) {
	as := assert.New(t)

	p := parse.String("hello").Fail("explode!")
	s, f := p.Parse("hello")
	as.Nil(s)
	as.NotNil(f)
	as.EqualError(f.Error, "explode!")
}

func TestParserSatisfy(t *testing.T) {
	as := assert.New(t)

	p := parse.String("hello").Satisfy(func(i parse.Input) (int, error) {
		if i[0] == '!' {
			return 1, nil
		}
		return 0, errors.New("mismatch")
	})

	s, f := p.Parse("hello!")
	as.NotNil(s)
	as.Nil(f)
	as.Equal(parse.Input("!"), s.Result)
}
