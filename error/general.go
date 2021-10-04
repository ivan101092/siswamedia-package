package error

var (
	generalBase = 10000
	generalErrs = map[int]SystemError{
		generalBase + 1: newError(generalBase+1, "Email tidak valid", "Invalid email"),
		generalBase + 2: newError(generalBase+2, "Password tidak valid", "Invalid password"),
		generalBase + 3: newError(generalBase+3, "Email tidak valid", "Invalid email"),
	}
)
