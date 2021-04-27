package parse

import "fmt"

type (
	// Input represents a Parser's input
	Input string

	arg = interface{}
)

// Error messages
const (
	ErrWrappedExpectation = "%s, got %s"
)

const (
	maxExpectedGot = 16
)

func (i Input) succeedWith(r Result) (*Success, *Failure) {
	return &Success{
		Result:    r,
		Remaining: i,
	}, nil
}

func (i Input) errExpected(msg string, args ...arg) error {
	got := i
	if len(got) > maxExpectedGot {
		got = got[0:maxExpectedGot] + "..."
	}
	errMsg := fmt.Sprintf(msg, args...)
	return fmt.Errorf(ErrWrappedExpectation, errMsg, got)
}

func (i Input) failMessage(msg string, args ...arg) (*Success, *Failure) {
	return i.failWith(fmt.Errorf(msg, args...))
}

func (i Input) failExpected(msg string, args ...arg) (*Success, *Failure) {
	return i.failWith(i.errExpected(msg, args...))
}

func (i Input) failThrough(f *Failure) (*Success, *Failure) {
	return i.failWith(f.Error)
}

func (i Input) failWith(err error) (*Success, *Failure) {
	return nil, &Failure{
		Error: err,
		Input: i,
	}
}
