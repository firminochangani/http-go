package client

import (
	"fmt"
	"io"
	"net"
	"net/url"
	"time"
)

type Client struct {
	Timeout time.Duration
}

type Response struct {
	Body io.ReadCloser
}

func (c *Client) validateDefaults() {
	if c.Timeout == 0 {
		c.Timeout = time.Second * 15
	}
}

func (c *Client) Get(path string) (*Response, error) {
	u, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	port := "80"
	if u.Port() != "" {
		port = u.Port()
	}

	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%s", u.Hostname(), port), c.Timeout)
	if err != nil {
		return nil, fmt.Errorf("unable to establish a connection to %s: %v", u.Hostname(), err)
	}

	_, err = conn.Write([]byte(fmt.Sprintf("GET %s", u.Path)))
	if err != nil {
		return nil, fmt.Errorf("unable to send a request to %s", u.Hostname())
	}

	return &Response{
		Body: &body{
			reader: conn,
			closer: conn,
		},
	}, nil
}

type body struct {
	reader io.Reader
	closer io.Closer
}

func (b *body) Read(p []byte) (n int, err error) {
	return b.reader.Read(p)
}

func (b *body) Close() error {
	return b.closer.Close()
}
