package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	setupE2E            = sync.Once{}
	stdOut              *bytes.Buffer
	stdErr              *bytes.Buffer
	stop                func()
	testUpstreamHandler = map[string]http.HandlerFunc{}
)

// ensureE2ESetup ensures that the setup for the end-to-end tests is done only once.
func ensureE2ESetup(t *testing.T) {
	setupE2E.Do(func() {
		// Setup the test upstream server.
		l, err := net.Listen("tcp", "127.0.0.1:8199")
		require.NoError(t, err)
		testUpstream := &httptest.Server{
			Listener: l,
			Config: &http.Server{Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				value := r.Header.Get("go-sdk-test-case")
				hander, ok := testUpstreamHandler[value]
				if !ok {
					log.Printf("testUpstreamHandler not found for %s", value)
					w.WriteHeader(http.StatusNotFound)
					return
				}
				hander(w, r)
			})},
		}
		testUpstream.Start()

		// Check if `envoy` is installed.
		_, err = exec.LookPath("envoy")
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
		stop = func() {
			testUpstream.Close()
			require.NoError(t, cmd.Process.Signal(os.Interrupt))
		}
		defer fmt.Println(stdOut.String())
		defer fmt.Println(stdErr.String())
	})
}

func TestHeaders(t *testing.T) {
	ensureE2ESetup(t)

	testUpstreamHandler["headers"] = func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "yes", r.Header.Get("foo"))
		require.Equal(t, "", r.Header.Get("multiple-values"))
		require.Equal(t, []string{"single"}, r.Header.Values("multiple-values-to-be-single"))

		w.Header().Set("this-is", "response-header")
		w.Header().Add("this-is-2", "A")
		w.Header().Add("this-is-2", "B")
		w.Header().Set("multiple-values2-to-be-single", "A")
		w.Header().Add("multiple-values2-to-be-single", "B")
		w.WriteHeader(http.StatusOK)
	}

	require.Eventually(t, func() bool {
		req, err := http.NewRequest("GET", "http://localhost:15002", bytes.NewBufferString("hello"))
		if err != nil {
			return false
		}
		req.Header.Set("go-sdk-test-case", "headers")
		req.Header.Set("foo", "value")
		req.Header.Add("multiple-values", "1234")
		req.Header.Add("multiple-values", "next")

		req.Header.Add("multiple-values2-to-be-single", "A")
		req.Header.Add("multiple-values2-to-be-single", "B")

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return false
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			return false
		}

		// Check if the response headers are as expected.
		if res.Header.Get("this-is") != "response-header" {
			fmt.Println("this-is:", res.Header.Values("this-is"))
			return false
		}
		if res.Header.Values("this-is-2") != nil {
			fmt.Println("this-is-2:", res.Header.Values("this-is-2"))
			return false
		}
		if toBeSingle := res.Header.Values("multiple-values-res-to-be-single"); len(toBeSingle) != 1 || toBeSingle[0] != "single" {
			fmt.Println("multiple-values-to-be-single:", toBeSingle)
			return false
		}
		return true
	}, 10*time.Second, 2*time.Second, "Envoy has not started: %s", stdOut.String())

	// Check if the log contains the expected output.
	requireEventuallyContainsMessages(t, stdOut,
		"foo: value",
		"multiple-values: 1234",
		"multiple-values: next",
		"this-is: response-header",
		"this-is-2: A",
		"this-is-2: B",
	)
}

func TestDelayFilter(t *testing.T) {
	ensureE2ESetup(t)

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
	ensureE2ESetup(t)

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

func TestBodies(t *testing.T) {
	ensureE2ESetup(t)

	testUpstreamHandler["bodies"] = func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		require.Equal(t, "XXXXXXXXXX", string(body)) // Request body should be replaced with 'X'.

		w.Header().Set("Content-Type", "text/plain")
		_, err = w.Write([]byte("example body\n"))
		require.NoError(t, err)
		w.WriteHeader(http.StatusOK)
	}

	// Make a request to the envoy proxy.
	require.Eventually(t, func() bool {
		req, err := http.NewRequest("GET", "http://localhost:15003", bytes.NewBufferString("0123456789"))
		if err != nil {
			return false
		}
		req.Header.Set("go-sdk-test-case", "bodies")
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return false
		}
		defer res.Body.Close()

		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			return false
		}
		require.Equal(t, "YYYYYYYYYYYYY", string(resBody))
		return res.StatusCode == http.StatusOK
	}, 5*time.Second, 100*time.Millisecond, "Envoy has not started: %s", stdOut.String())

	// Check if the log contains the expected output.
	requireEventuallyContainsMessages(t, stdOut,
		"entire request body: 0123456789",
		"request body read 2 bytes offset at 0: \"01\"",
		"request body read 2 bytes offset at 2: \"23\"",
		"request body read 2 bytes offset at 4: \"45\"",
		"request body read 2 bytes offset at 6: \"67\"",
		"request body read 2 bytes offset at 8: \"89\"",

		"entire response body: example body",
		"response body read 2 bytes offset at 0: \"ex\"",
		"response body read 2 bytes offset at 2: \"am\"",
		"response body read 2 bytes offset at 4: \"pl\"",
		"response body read 2 bytes offset at 6: \"e \"",
		"response body read 2 bytes offset at 8: \"bo\"",
		"response body read 2 bytes offset at 10: \"dy\"",
	)
}

func requireEventuallyContainsMessages(t *testing.T, buf *bytes.Buffer, messages ...string) {
	for _, msg := range messages {
		require.Eventually(t, func() bool {
			return strings.Contains(buf.String(), msg)
		}, 3*time.Second, 100*time.Millisecond, "Expected message \"%s\" not found in buffer\n%s", msg, buf.String())
	}
}
