package helper

import (
	"RestGoTest/src/constant"

	"net/http"
)

var StatusCodeMapping = map[string]int{

	// OTP
	constant.OptExists:   409,
	constant.OtpUsed:     409,
	constant.OtpNotValid: 400,

	// User
	constant.EmailExists:      409,
	constant.UsernameExists:   409,
	constant.RecordNotFound:   404,
	constant.PermissionDenied: 403,
}

func TranslateErrorToStatusCode(err error) int {
	value, ok := StatusCodeMapping[err.Error()]
	if !ok {
		return http.StatusInternalServerError
	}
	return value
}
