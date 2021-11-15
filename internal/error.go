package internal

import "github.com/cockroachdb/errors"

type Error struct {
	message string
	inner   error
	outer   error
}

func NewError(message string) func() *Error {
	return func() *Error {
		return &Error{
			message: message,
			outer:   errors.NewWithDepth(1, message),
		}
	}
}

func (self *Error) Wrap(err error) *Error {
	self.inner = err
	self.outer = errors.WrapWithDepth(1, self.inner, self.message)
	return self
}

func (self *Error) WrapWithDepth(depth int, err error) *Error {
	self.inner = err
	self.outer = errors.WrapWithDepth(depth+1, self.inner, self.message)
	return self
}

func (self Error) Outer() error {
	return self.outer
}

func (self Error) Inner() error {
	return self.inner
}

func (self Error) Error() string {
	return self.Outer().Error()
}

func (self Error) Unwrap() error {
	return self.Inner()
}

func (self Error) Is(err error) bool {
	if outer, ok := err.(*Error); ok {
		return self.message == outer.message
	}

	return false
}
