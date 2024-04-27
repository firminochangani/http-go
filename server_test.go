package http_go_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	http_go "github.com/flowck/http-go"
)

func TestServer_ListenAndServe(t *testing.T) {
	t.Run("server_is_initialised_with_the_default_values", func(t *testing.T) {
		srv := http_go.Server{}
		go func() { _ = srv.ListenAndServe() }()
		time.Sleep(time.Millisecond * 250)

		assert.NotNil(t, srv.Ctx)
		assert.NotEmpty(t, srv.Addr, "a port is automatically chosen")

		require.NoError(t, srv.Shutdown())
	})
}
