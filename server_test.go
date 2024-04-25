package http_go_test

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	http_go "github.com/flowck/http-go"
)

func TestServer_GET(t *testing.T) {
	router := http_go.NewServerDefaultNaiveRouter()
	router.GET("/people", func(r *http_go.Request, w *http_go.Response) error {
		return w.Write([]byte("Hello World"))
	})

	srv := http_go.Server{
		Addr:   ":6060",
		Router: router,
	}
	client := http.Client{
		Timeout: time.Millisecond * 500,
	}
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		srvErr := srv.ListenAndServe()
		if srvErr != nil && !errors.Is(srvErr, http_go.ErrServerIsClosed) {
			t.Log(srvErr)
		}
	}(wg)

	resp, err := client.Get("http://localhost:6060/people")
	require.NoError(t, err)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.NoError(t, resp.Body.Close())

	assert.Equal(t, http_go.StatusCodeOK.Code, resp.StatusCode)
	assert.Equal(t, "Hello World", string(body))

	require.NoError(t, srv.Shutdown())

	_, err = client.Get("http://localhost:6060/people")
	require.Error(t, err, "no request shall be handled after .Shutdown() gets called")
	wg.Wait()
}

func TestReqParser(t *testing.T) {
	message := string(getReq())

	t.Log(len(strings.Split(message, "\n")))
}

func getReq() []byte {
	return []byte(`
POST / HTTP/1.1
Host: www.example.re
Content-Type: multipart/form-data; boundary=”NextField”
Content-Length: 125

--NextField
Content-Disposition: form-data; name=”Job”

100

-- NextField
Content-Disposition: form-data; name=”Priority”

2
`)
}
