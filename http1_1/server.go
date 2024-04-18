package http1_1

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/url"
	"strings"
	"time"
)

const (
	MethodGET     = "GET"
	MethodHEAD    = "HEAD"
	MethodPOST    = "POST"
	MethodPUT     = "PUT"
	MethodDELETE  = "DELETE"
	MethodCONNECT = "CONNECT"
	MethodOPTIONS = "OPTIONS"
	MethodTRACE   = "TRACE"
	MethodPATCH   = "PATCH"
)

type Header map[string]string

func (h Header) Set(name, value string) {
	h[name] = value
}

type Request struct {
	Method  string
	Headers Header
	URL     *url.URL
}

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

		_, err := r.conn.Write([]byte(fmt.Sprintf("HTTP/1.1 %d %s\n%s\n\n", r.statusCode.Code, r.statusCode.Name, headers)))
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

type Server struct {
	Addr   string
	Router Router

	ctx      context.Context
	listener net.Listener
}

func (s *Server) ListenAndServe() error {
	var err error
	s.listener, err = net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}

	if s.ctx == nil {
		s.ctx = context.Background()
	}

	s.acceptLoop()

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.listener.Close()
}

func (s *Server) acceptLoop() {
	for {
		conn, err := s.listener.Accept()
		if errors.Is(err, net.ErrClosed) {
			break
		}

		go s.handleRequest(context.WithoutCancel(s.ctx), conn)
	}
}

func (s *Server) handleRequest(ctx context.Context, conn net.Conn) {
	message := make([]byte, 1024*2)
	r := &Request{}
	w := &Response{
		conn:    conn,
		Headers: Header{},
	}

	//TODO: assert whether in a post request with multipart payload
	// the headers are read all at once or read line by line
	n, err := conn.Read(message)
	if err != nil {
		log.Println("unable to read the request", message)
		s.closeConn(conn)
		return
	}

	parseRequest(r, message[:n])
	err = s.Router.Handle(r, w)
	if err != nil {
		fmt.Println(err)
	}

	s.closeConn(conn)
}

func parseRequest(r *Request, message []byte) *Request {
	r.Headers = make(Header)

	line := ""
	var lineHeader []string
	lineCount := 0
	for i := 0; i < len(message); i++ {
		if message[i] == 10 {
			// request's first line
			if lineCount == 0 {
				r.Method = strings.TrimSpace(strings.Split(line[:7], " ")[0])
				u, err := url.Parse(strings.Split(line, " ")[1])
				if err != nil {
					log.Println("unable to parse url: ", err)
				}

				r.URL = u
			} else {
				lineHeader = strings.SplitN(line, ":", 2)
				if len(lineHeader) > 1 {
					r.Headers[lineHeader[0]] = lineHeader[1]
				} else {
					r.Headers[lineHeader[0]] = ""
				}
			}
			line = ""
			lineCount++
		} else {
			line += string(message[i])
		}
	}

	return r
}

func (s *Server) closeConn(conn net.Conn) {
	err := conn.Close()
	if err != nil {
		fmt.Println("unable to close the connection successfully", err)
	}
}
