package http_go

type StatusCode struct {
	Name string
	Code int
}

var (
	StatusCodeOK              = StatusCode{Name: "OK", Code: 200}
	StatusCodeNotFound        = StatusCode{Name: "Not Found", Code: 404}
	StatusInternalServerError = StatusCode{Name: "Internal Server Error", Code: 500}
)

func newStatusCode(code int) StatusCode {
	switch code {
	case StatusCodeOK.Code:
		return StatusCodeOK
	case StatusCodeNotFound.Code:
		return StatusCodeNotFound
	case StatusInternalServerError.Code:
		return StatusInternalServerError
	default:
		return StatusCodeOK
	}
}
