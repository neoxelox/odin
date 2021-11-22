package internal

import "github.com/cockroachdb/errors"

type Error struct {
	message string
	inner   error
}

func NewError(message string) func() *Error {
	return func() *Error {
		return &Error{
			message: message,
			inner:   errors.NewWithDepth(1, message),
		}
	}
}

func (self *Error) Wrap(err error) *Error {
	self.inner = errors.WrapWithDepth(1, err, self.message)
	return self
}

func (self *Error) WrapWithDepth(depth int, err error) *Error {
	self.inner = errors.WrapWithDepth(depth+1, err, self.message)
	return self
}

func (self *Error) As(err error) *Error {
	if other, ok := err.(*Error); ok {
		self.message = other.message
	} else {
		self.message = err.Error()
	}
	self.inner = errors.WrapWithDepth(1, err, self.message)
	return self
}

func (self *Error) AsWithDepth(depth int, err error) *Error {
	if other, ok := err.(*Error); ok {
		self.message = other.message
	} else {
		self.message = err.Error()
	}
	self.inner = errors.WrapWithDepth(depth+1, err, self.message)
	return self
}

func (self *Error) With(message string) *Error {
	self.inner = errors.NewWithDepth(1, self.message+": "+message)
	return self
}

func (self Error) Error() string {
	return self.Unwrap().Error()
}

func (self Error) Unwrap() error {
	return self.inner
}

func (self Error) Is(err error) bool {
	if other, ok := err.(*Error); ok {
		return self.message == other.message
	}

	return false
}
