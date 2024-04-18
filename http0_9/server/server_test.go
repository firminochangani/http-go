package server_test

import (
	"errors"
	"io"
	"log/slog"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"http_v0_9/client"
	"http_v0_9/server"
)

func TestServer_Lifecycle(t *testing.T) {
	wg := &sync.WaitGroup{}
	cli := &client.Client{}
	srv, srvErr := server.NewServer("localhost", 8090, "../static", slog.New(slog.NewTextHandler(os.Stdout, nil)))
	require.NoError(t, srvErr)

	wg.Add(1)
	go serverRunner(t, wg, srv)

	t.Run("get_the_home_page", func(t *testing.T) {
		resp, err := cli.Get("http://localhost:8090")
		require.NoError(t, err)

		body := readBody(t, resp)
		assert.NotEmpty(t, body)
	})

	t.Run("get_non_existent_page", func(t *testing.T) {
		resp, err := cli.Get("http://localhost:8090/some/cool/page")
		require.NoError(t, err)

		body := readBody(t, resp)
		assert.Containsf(t, body, server.ErrNotFound.Error(), "a 404 page gets returned")
	})

	// Clean up
	require.NoError(t, srv.Shutdown())
	wg.Wait()
}

func readBody(t *testing.T, response *client.Response) string {
	require.NotNil(t, response)

	b, err := io.ReadAll(response.Body)
	require.NoError(t, err)

	return string(b)
}

func serverRunner(t *testing.T, wg *sync.WaitGroup, srv *server.Server) {
	defer wg.Done()

	t.Log("the server is running")

	err := srv.Listen()
	if err != nil && !errors.Is(err, server.ErrHttpServerClosed) {
		t.Logf("the server stopped with the following error: %v", err)
		return
	}

	t.Log("the server got shutdown")
}
