package http_go

import "fmt"

type StatusCode struct {
	Name string
	Code int
}

func (c StatusCode) String() string {
	return fmt.Sprintf("%d %s", c.Code, c.Name)
}

var (
	StatusCodeOK              = StatusCode{Name: "OK", Code: 200}
	StatusBadRequest          = StatusCode{Name: "Bad Request", Code: 400}
	StatusCodeNotFound        = StatusCode{Name: "Not Found", Code: 404}
	StatusInternalServerError = StatusCode{Name: "Internal Server Error", Code: 500}
	StatusAccepted            = StatusCode{Name: "Accepted", Code: 202}
)

func newStatusCode(code int) StatusCode {
	switch code {
	// 2xx
	case StatusCodeOK.Code:
		return StatusCodeOK
	case StatusAccepted.Code:
		return StatusAccepted
	// 4xx
	case StatusBadRequest.Code:
		return StatusBadRequest
	case StatusCodeNotFound.Code:
		return StatusCodeNotFound
	// 5xx
	case StatusInternalServerError.Code:
		return StatusInternalServerError
	default:
		return StatusCodeOK
	}
}
