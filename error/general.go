package error

import "errors"

var (
	generalBase = 10000
	generalErrs = map[int]SystemError{
		generalBase + 1: newError(generalBase+1, errors.New("invalid_email"), "Email tidak valid", "Invalid email"),
		generalBase + 2: newError(generalBase+2, errors.New("invalid_password"), "Password tidak valid", "Invalid password"),
		generalBase + 3: newError(generalBase+3, errors.New("invalid_email"), "Email tidak valid", "Invalid email"),
	}
)
