package otp

import "github.com/neoxelox/odin/internal"

const (
	OTP_MESSAGE = `Tu código de verificación para Community es: %s. Por favor, nunca compartas este código con nadie.`
)

var (
	ErrGeneric     = internal.NewError("OTP execution failed")
	ErrInvalid     = internal.NewError("OTP is invalid")
	ErrAlreadySent = internal.NewError("OTP recently sent")
	ErrMaxAttempts = internal.NewError("Maximum OTP attempts reached")
	ErrWrongCode   = internal.NewError("OTP wrong code")
)
