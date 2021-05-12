package parse_test

import (
	"fmt"
	"testing"

	"github.com/caravan/kombi/parse"
	"github.com/stretchr/testify/assert"
)

type Assertions struct {
	*assert.Assertions
}

func NewAssert(t *testing.T) *Assertions {
	return &Assertions{
		Assertions: assert.New(t),
	}
}

func (as *Assertions) Success(s *parse.Success, f *parse.Failure) {
	as.NotNil(s)
	as.Nil(f)
}

func (as *Assertions) SuccessResult(
	s *parse.Success, f *parse.Failure, r interface{},
) {
	as.Success(s, f)
	as.Equal(r, s.Result)
}

func (as *Assertions) SuccessResults(
	s *parse.Success, f *parse.Failure, r ...interface{},
) {
	as.Success(s, f)
	as.Equal(len(r), len(s.Result.(parse.Results)))
	for i, e := range r {
		as.Equal(e, s.Result.(parse.Results)[i])
	}
}

func (as *Assertions) Failure(s *parse.Success, f *parse.Failure) {
	as.Nil(s)
	as.NotNil(f)
}

func (as *Assertions) FailureWrapped(
	s *parse.Success, f *parse.Failure, msg string, args ...interface{},
) {
	as.Failure(s, f)
	as.Wrapped(f.Error, msg, args...)
}

func (as *Assertions) FailureError(
	s *parse.Success, f *parse.Failure, msg string, args ...interface{},
) {
	as.Failure(s, f)
	as.EqualError(f.Error, msg, args...)
}

func (as *Assertions) Wrapped(err error, msg string, args ...interface{}) {
	allArgs := append([]interface{}{msg}, args...)
	as.EqualError(err, fmt.Sprintf(parse.ErrWrappedExpectation, allArgs...))
}
