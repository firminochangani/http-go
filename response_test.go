package http_go

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResponse_Write(t *testing.T) {
	mockNetConn := &MockNetConn{}
	w := Response{
		Headers: make(Header),
		conn:    mockNetConn,
	}
	err := w.Write([]byte("Hello world"))
	require.NoError(t, err)

	assert.Contains(t, mockNetConn.messageWritten, "Hello world")
}

type MockNetConn struct {
	messageWritten string
}

func (m *MockNetConn) Read(b []byte) (n int, err error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockNetConn) Write(b []byte) (n int, err error) {
	m.messageWritten += string(b)
	return len(b), nil
}

func (m *MockNetConn) Close() error {
	//TODO implement me
	panic("implement me")
}

func (m *MockNetConn) LocalAddr() net.Addr {
	//TODO implement me
	panic("implement me")
}

func (m *MockNetConn) RemoteAddr() net.Addr {
	//TODO implement me
	panic("implement me")
}

func (m *MockNetConn) SetDeadline(t time.Time) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockNetConn) SetReadDeadline(t time.Time) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockNetConn) SetWriteDeadline(t time.Time) error {
	//TODO implement me
	panic("implement me")
}
