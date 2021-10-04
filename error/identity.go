package error

import "errors"

var (
	identityBase = 11000
	identityErrs = map[int]SystemError{
		identityBase + 1: newError(identityBase+1, errors.New("invalid_credential"), "Email/password tidak valid", "Invalid email/password"),
	}
)
