package http

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
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
	Addr   string
	Router Router

	isRunning atomic.Bool
	listener  net.Listener
	done      chan interface{}
	Ctx       context.Context
}

func (s *Server) setServerDefaults() {
	if s.Ctx == nil {
		s.Ctx = context.Background()
	}

	s.isRunning = atomic.Bool{}
	s.isRunning.Store(true)
	s.done = make(chan interface{})
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

func (s *Server) Shutdown(ctx context.Context) error {
	if s.listener == nil {
		return nil
	}

	if s.isRunning.Load() {
		close(s.done)
		s.isRunning.Store(false)
	}

	return s.listener.Close()
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

			go s.handleRequest(context.WithoutCancel(s.Ctx), conn)
		}
	}
}

func (s *Server) handleRequest(ctx context.Context, conn net.Conn) {
	message := make([]byte, 1024*2)
	r := &Request{
		Context: ctx,
	}

	w := &Response{
		conn: conn,
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

func (s *Server) closeConn(conn net.Conn) {
	err := conn.Close()
	if err != nil {
		fmt.Println("unable to close the connection successfully", err)
	}
}
