package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	setupEnvoy = sync.Once{}
	stdOut     *bytes.Buffer
	stdErr     *bytes.Buffer
	stop       func()
)

func ensureEnvoy(t *testing.T) {
	setupEnvoy.Do(func() {
		// Check if `envoy` is installed.
		_, err := exec.LookPath("envoy")
		require.NoError(t, err, "envoy binary not found. Please install it from containers at https://github.com/envoyproxyx/envoyx/pkgs/container/envoy")

		// Check if a binary named main exists.
		_, err = os.Stat("./main")
		require.NoError(t, err, "./main not found. Please build it.")

		// Check if envoy.yaml exists.
		_, err = os.Stat("./envoy.yaml")
		require.NoError(t, err, "./envoy.yaml not found. Please create it.")

		cmd := exec.Command("envoy", "--concurrency", "1", "-c", "./envoy.yaml")
		stdOut, stdErr = new(bytes.Buffer), new(bytes.Buffer)
		cmd.Stdout = stdOut
		cmd.Stderr = stdErr
		require.NoError(t, cmd.Start())
		stop = func() { require.NoError(t, cmd.Process.Signal(os.Interrupt)) }
		defer fmt.Println(stdOut.String())
		defer fmt.Println(stdErr.String())
	})
}

func TestDelayFilter(t *testing.T) {
	ensureEnvoy(t)

	// Make four requests to the envoy proxy.
	wg := new(sync.WaitGroup)
	wg.Add(4)
	for i := 0; i < 4; i++ {
		go func() {
			defer wg.Done()
			require.Eventually(t, func() bool {
				req, err := http.NewRequest("GET", "http://localhost:15001", bytes.NewBufferString("hello"))
				if err != nil {
					return false
				}
				res, err := http.DefaultClient.Do(req)
				if err != nil {
					return false
				}
				defer res.Body.Close()
				return res.StatusCode == http.StatusOK
			}, 10*time.Second, 2*time.Second, "Envoy has not started: %s", stdOut.String())
		}()
	}
	wg.Wait()

	// Check if the log contains the expected output.
	requireEventuallyContainsMessages(t, stdOut,
		"EventHttpRequestHeaders returning StopAllIterationAndBuffer with id 1",
		"blocking for 1 second at EventHttpRequestHeaders with id 1",
		"calling ContinueRequest with id 1",
		"EventHttpRequestBody called with id 1",
		"EventHttpResponseHeaders called with id 1",
		"EventHttpRequestHeaders called with id 2",
		"EventHttpRequestBody returning StopIterationAndBuffer with id 2",
		"blocking for 1 second at EventHttpRequestBody with id 2",
		"calling ContinueRequest with id 2",
		"EventHttpResponseBody called with id 2",
		"EventHttpRequestHeaders called with id 3",
		"EventHttpRequestBody called with id 3",
		"EventHttpResponseHeaders returning StopAllIterationAndBuffer with id 3",
		"blocking for 1 second at EventHttpResponseHeaders with id 3",
		"calling ContinueResponse with id 3",
		"EventHttpResponseBody called with id 3",
		"EventHttpRequestHeaders called with id 4",
		"EventHttpRequestBody called with id 4",
		"EventHttpResponseHeaders called with id 4",
		"blocking for 1 second at EventHttpResponseBody with id 4",
		"EventHttpResponseBody returning StopIterationAndBuffer with id 4",
		"calling ContinueResponse with id 4",
	)
}

func TestHelloWorld(t *testing.T) {
	ensureEnvoy(t)

	// Make a request to the envoy proxy.
	require.Eventually(t, func() bool {
		req, err := http.NewRequest("GET", "http://localhost:15000", bytes.NewBufferString("hello"))
		if err != nil {
			return false
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return false
		}
		defer res.Body.Close()
		return res.StatusCode == http.StatusOK
	}, 5*time.Second, 100*time.Millisecond, "Envoy has not started: %s", stdOut.String())

	// Check if the log contains the expected output.
	requireEventuallyContainsMessages(t, stdOut,
		"helloWorldHttpFilter.NewHttpFilterInstance called",
		"helloWorldHttpFilterInstance.EventHttpRequestHeaders called",
		"helloWorldHttpFilterInstance.EventHttpRequestBody called",
		"helloWorldHttpFilterInstance.EventHttpResponseHeaders called",
		"helloWorldHttpFilterInstance.EventHttpResponseBody called",
		"helloWorldHttpFilterInstance.EventHttpDestroy called",
	)
}

func requireEventuallyContainsMessages(t *testing.T, buf *bytes.Buffer, messages ...string) {
	for _, msg := range messages {
		require.Eventually(t, func() bool {
			return strings.Contains(buf.String(), msg)
		}, 3*time.Second, 100*time.Millisecond, "Expected message \"%s\" not found in buffer\n%s", msg, buf.String())
	}
}
