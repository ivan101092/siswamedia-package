package error

type SystemError struct {
	Error   error              `json:"error"` // Original errors from another services
	Code    int                `json:"code"`
	Message SystemErrorMessage `json:"message"`
}

type SystemErrorMessage struct {
	ID string `json:"id"`
	En string `json:"error"`
}

var (
	baseErrs = map[int]SystemError{
		0: newError(0, "Terjadi kendala pada server", "Internal server error"),
	}
)

func newError(code int, messageID, messageEn string) SystemError {
	return SystemError{
		Code: code,
		Message: SystemErrorMessage{
			ID: messageID,
			En: messageEn,
		},
	}
}

func GetError(code int) SystemError {
	if code < 11000 {
		return generalErrs[code]
	} else if code < 12000 {
		return identityErrs[code]
	}

	return baseErrs[0]
}
