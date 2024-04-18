package client_test

import (
	"errors"
	"fmt"
	"io"
	"net"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"http_v0_9/client"
)

func TestClient_Get(t *testing.T) {
	expectedResponse := gofakeit.Sentence(100)
	stopSrv, srvAddr := mockTcpServer(t, expectedResponse, 0)
	defer func() { require.NoError(t, stopSrv()) }()

	c := client.Client{
		Timeout: time.Second * 2,
	}
	resp, err := c.Get(fmt.Sprintf("http://%s/contacts/phones.html", srvAddr))
	require.NoError(t, err)
	assert.NotNil(t, resp)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, expectedResponse, string(body))
}

func TestClient_Get_RealServer(t *testing.T) {
	c := client.Client{
		Timeout: time.Second * 2,
	}
	resp, err := c.Get("http://localhost:8080/index.html")
	require.NoError(t, err)
	assert.NotNil(t, resp)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.NotEmpty(t, body)
}

func mockTcpServer(t *testing.T, response string, timeout time.Duration) (func() error, string) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	go func() {
		for {
			conn, err := l.Accept()
			if errors.Is(err, net.ErrClosed) {
				t.Log("server has been closed")
				break
			}

			if err != nil {
				t.Logf("unable to accept connection: %v", err)
				continue
			}

			if timeout > 0 {
				t.Logf("timeout per request has been set to %s. Waiting...", timeout)
				time.Sleep(timeout)
			}

			reqBuff := make([]byte, 1024)
			_, err = conn.Read(reqBuff)
			if err != nil {
				t.Logf("unable to read the request payload: %v", err)
				return
			}

			n, err := conn.Write([]byte(response))
			if err != nil {
				t.Logf("unable to write to the client: %v", err)
			}

			t.Logf("response written to the client: %d bytes sent", n)

			err = conn.Close()
			if err != nil {
				t.Logf("an error occurred while closing the connection: %v", err)
			}
		}
	}()

	return l.Close, l.Addr().String()
}
