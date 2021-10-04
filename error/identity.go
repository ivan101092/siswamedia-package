package error

var (
	identityBase = 11000
	identityErrs = map[int]SystemError{
		identityBase + 1: newError(identityBase+1, "Email/password tidak valid", "Invalid email/password"),
	}
)
