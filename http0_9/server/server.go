package server

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"strings"
)

const (
	// GET + <blank_space> + URL
	maxBytesPerRequest = 3 + 1 + (1024 * 2)
)

var (
	ErrHttpServerClosed    = errors.New("the http server has been closed")
	ErrMethodNotAllowed    = errors.New("405 - Method Not Allowed")
	ErrUriTooLong          = errors.New("414 - URI Too Long")
	ErrInternalServerError = errors.New("500 - Internal Server Error")
	ErrNotFound            = errors.New("404 - Not Found")
)

type Server struct {
	rootDir  string
	logger   *slog.Logger
	listener net.Listener
}

type Request struct {
	path   string
	method string
}

func NewServer(host string, port int, rootDir string, logger *slog.Logger) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return nil, fmt.Errorf("unable to listen to %s:%d", host, port)
	}

	return &Server{
		logger:   logger,
		rootDir:  rootDir,
		listener: listener,
	}, nil
}

func (s *Server) Listen() error {
	for {
		conn, err := s.listener.Accept()
		if errors.Is(err, net.ErrClosed) {
			s.logger.Debug("server closed and it's no longer accepting connection")
			return ErrHttpServerClosed
		}

		if err != nil {
			s.logger.Error("unable to accept connection", "error", err)
			continue
		}

		go func() {
			err := s.handleConn(conn)
			if err != nil {
				s.logger.Error("error while handling connection: ", "error: ", err)
			}
		}()
	}
}

func (s *Server) Shutdown() error {
	return s.listener.Close()
}

func (s *Server) handleConn(conn net.Conn) error {
	readBuff := make([]byte, maxBytesPerRequest)
	n, err := conn.Read(readBuff)
	if err != nil {
		return sendResponse(conn, []byte("unable to read request payload"))
	}

	// parse request body
	method := string(readBuff[:3])
	if strings.ToUpper(method) != "GET" {
		return sendResponse(conn, []byte(ErrMethodNotAllowed.Error()))
	}

	path := strings.TrimSpace(string(readBuff[4:n]))
	if len(path) > (1024 * 2) {
		return s.sendResponseAndCloseConn(conn, []byte(ErrUriTooLong.Error()))
	}

	if path == "" {
		path = "/index.html"
	}

	err = s.readFileInStreams(path, func(chunks []byte) error {
		return sendResponse(conn, chunks)
	})
	if errors.Is(err, os.ErrNotExist) {
		return s.sendResponseAndCloseConn(conn, []byte(ErrNotFound.Error()))
	}
	if err != nil {
		return s.sendResponseAndCloseConn(conn, []byte(ErrInternalServerError.Error()))
	}

	s.logger.Info("request processed", "method", method, "path", path, "file", fmt.Sprintf("%s%s", s.rootDir, path))

	return s.closeConn(conn)
}

func (s *Server) readFileInStreams(path string, handler func(chunks []byte) error) error {
	file, err := os.Open(fmt.Sprintf("%s%s", s.rootDir, path))
	if err != nil {
		return err
	}

	defer func() { _ = file.Close() }()

	var chunks []byte
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		chunks = scanner.Bytes()
		if scanner.Err() != nil && errors.Is(scanner.Err(), io.EOF) {
			break
		}

		if scanner.Err() != nil {
			return scanner.Err()
		}

		err = handler(chunks)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) sendResponseAndCloseConn(conn net.Conn, body []byte) error {
	err := sendResponse(conn, body)
	if err != nil {
		return s.closeConn(conn)
	}

	return s.closeConn(conn)
}

func (s *Server) closeConn(conn net.Conn) error {
	err := conn.Close()
	if err != nil {
		return fmt.Errorf("unable to close connection: %v", err)
	}

	return nil
}

func sendResponse(conn net.Conn, body []byte) error {
	_, err := conn.Write(body)
	if err != nil {
		return errors.New("unable to send response")
	}

	return nil
}
