package csv

import "fmt"

type err string

func (e err) Error() string {
	return string(e)
}

const (
	ErrContract = err("invalid contract")
)

// NewError is just a syntax suggar for fmt.Error("%w: %s", err, cause)
// Used to give more context to the err error.
func NewError(err error, cause string) error {
	return fmt.Errorf("%w: %s", err, cause)
}
