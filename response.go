package http_go

import (
	"fmt"
	"net"
	"time"
)

type Response struct {
	Headers Header

	responseWritten bool
	conn            net.Conn
	statusCode      StatusCode
}

func (r *Response) Write(message []byte) error {
	if !r.responseWritten && r.statusCode.Code == 0 {
		r.statusCode = StatusCodeOK
	}

	r.Headers.Set("Date", time.Now().String())

	// write the status code set previously if and only if no previous response has been set to the client
	if !r.responseWritten && r.statusCode.Code > 0 {
		headers := ""
		for name, value := range r.Headers {
			headers += fmt.Sprintf("%s: %s\n", name, value)
		}

		_, err := r.conn.Write([]byte(fmt.Sprintf("HTTP/1.1 %d %s\n%s\n", r.statusCode.Code, r.statusCode.Name, headers)))
		if err != nil {
			return err
		}
	}

	r.responseWritten = true
	_, err := r.conn.Write(message)
	if err != nil {
		return err
	}

	return nil
}

func (r *Response) WriteStatus(statusCode int) {
	r.statusCode = newStatusCode(statusCode)
}
