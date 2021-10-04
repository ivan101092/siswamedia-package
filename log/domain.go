package log

const (
	ProcessIDContextKey = "processID"
)

// Data Model for tracking incoming request
type RequestLogModel struct {
	ProcessID  string
	UserID     int
	IP         string
	Method     string
	URL        string
	ReqHeader  interface{}
	ReqBody    interface{}
	RespHeader interface{}
	RespBody   interface{}
	Error      interface{}
	StatusCode int
	Duration   int64
}

type ErrorData struct {
	Location string
	Error    string
}

type Errors []ErrorData

type Log struct {
	LogToTerminal bool
	Location      string
	FileFormat    string
	FileLinkName  string
	MaxAge        int
	RotationFile  int
	UseStackTrace bool
}
