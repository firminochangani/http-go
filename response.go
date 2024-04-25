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
		headers := "\r\n"
		for name, value := range r.Headers {
			headers += fmt.Sprintf("%s: %s\r\n", name, value)
		}
		headers += "\r\n"

		//TODO: extra allocation
		res := fmt.Sprintf("HTTP/1.1 %s%s%s", r.statusCode, headers, message)
		_, err := r.conn.Write([]byte(res))
		if err != nil {
			return err
		}
	} else {
		_, err := r.conn.Write(message)
		if err != nil {
			return err
		}
	}

	r.responseWritten = true

	return nil
}

func (r *Response) WriteStatus(statusCode int) {
	r.statusCode = newStatusCode(statusCode)
}
