package http_go

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/textproto"
	"net/url"
	"strings"
	"sync/atomic"
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

var (
	ErrServerIsClosed        = errors.New("server is closed")
	ErrServerContextIsClosed = errors.New("server's context is closed")
)

type Server struct {
	MaxHeaderBytes int
	Addr           string
	Router         Router
	Ctx            context.Context

	isRunning atomic.Bool
	listener  net.Listener
	done      chan interface{}
}

func (s *Server) ListenAndServe() error {
	var err error
	s.listener, err = net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}

	s.setServerDefaults()

	return s.acceptLoop()
}

func (s *Server) Shutdown() error {
	if s.listener == nil {
		return nil
	}

	if s.isRunning.Load() {
		close(s.done)
		s.isRunning.Store(false)
	}

	return s.listener.Close()
}

func (s *Server) setServerDefaults() {
	if s.Ctx == nil {
		s.Ctx = context.Background()
	}

	s.isRunning = atomic.Bool{}
	s.isRunning.Store(true)
	s.done = make(chan interface{})
}

func (s *Server) acceptLoop() error {
	for {
		select {
		case <-s.done:
			return ErrServerIsClosed
		case <-s.Ctx.Done():
			return ErrServerContextIsClosed
		default:
			conn, err := s.listener.Accept()
			if errors.Is(err, net.ErrClosed) {
				break
			}

			if err != nil {
				fmt.Println(err)
				continue
			}

			log.Println("handling request")
			go s.handleRequest(context.WithoutCancel(s.Ctx), conn)
		}
	}
}

func (s *Server) handleRequest(ctx context.Context, conn net.Conn) {
	r := &Request{
		ctx:     ctx,
		Headers: Header{},
	}

	w := &Response{
		conn:       conn,
		Headers:    Header{},
		statusCode: StatusCode{},
	}

	connReader := bufio.NewReader(conn)

	err := parseRequest(r, connReader)
	if err != nil {
		w.WriteStatus(StatusBadRequest.Code)
		_ = w.Write([]byte(err.Error()))
		s.closeConn(conn)
		return
	}

	err = s.Router.Handle(r, w)
	if err != nil && !w.responseWritten {
		_ = w.Write([]byte(fmt.Sprintf("Unhandled error: %s", err)))
	} else if !w.responseWritten {
		_ = w.Write([]byte(""))
	}

	s.closeConn(conn)
}

func (s *Server) closeConn(conn net.Conn) {
	err := conn.Close()
	if err != nil {
		fmt.Println("unable to close the connection successfully", err)
	}
}

func parseRequest(r *Request, connReader *bufio.Reader) error {
	tp := textproto.NewReader(connReader)

	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Messages#request_line
	requestLine, err := tp.ReadLine()
	if err != nil {
		panic(err)
	}

	requestLineParts := strings.Split(requestLine, " ")
	if len(requestLineParts) < 2 {
		return fmt.Errorf("error %s", StatusBadRequest)
	}

	// Set the method
	r.Method = strings.ToUpper(requestLineParts[0])

	// Parse the request url
	requestUrl, err := url.Parse(requestLineParts[1])
	if err != nil {
		return fmt.Errorf("error %s", StatusBadRequest)
	}
	r.URL = requestUrl

	// Set the protocols version. This could be useful if I ever intend to support HTTP/2
	if len(requestLineParts) >= 3 {
		r.Proto = strings.TrimSpace(requestLineParts[2])
	}

	headers, err := tp.ReadMIMEHeader()
	if err != nil {
		log.Println(err)
	} else {
		for k, v := range headers {
			r.Headers.Set(strings.ToLower(k), strings.Join(v, " "))
		}
	}

	if host, exists := r.Headers.Get("host"); exists {
		r.Host = host
	}

	return nil
}

func parseMessageToRequestBack(r *Request, message []byte) *Request {
	r.Headers = make(Header)

	line := ""
	lineCount := 0
	var lineHeader []string
	for i := 0; i < len(message); i++ {
		//nolint
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
