package otp

import "github.com/neoxelox/odin/internal"

const (
	OTP_MESSAGE = `Tu código de verificación es: %s. Por favor, nunca compartas este código con nadie.`
)

var (
	ErrGeneric     = internal.NewError("OTP execution failed")
	ErrAlreadySend = internal.NewError("OTP recently sent")
	ErrInvalidOTP  = internal.NewError("OTP is invalid")
	ErrMaxAttempts = internal.NewError("Maximum OTP attempts reached")
	ErrWrongCode   = internal.NewError("OTP wrong code")
)
