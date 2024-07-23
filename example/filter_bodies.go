package main

import (
	"fmt"
	"io"

	"github.com/envoyproxyx/go-sdk/envoy"
)

// bodiesHttpFilter implements envoy.HttpFilter.
//
// This is to demonstrate how to use body manipulation APIs.
type bodiesHttpFilter struct{}

func newbodiesHttpFilter(string) envoy.HttpFilter { return &bodiesHttpFilter{} }

// NewHttpFilterInstance implements envoy.HttpFilter.
func (f *bodiesHttpFilter) NewHttpFilterInstance(envoyFilter envoy.EnvoyFilterInstance) envoy.HttpFilterInstance {
	return &bodiesHttpFilterInstance{envoyFilter: envoyFilter}
}

// Destroy implements envoy.HttpFilter.
func (f *bodiesHttpFilter) Destroy() {}

// bodiesHttpFilterInstance implements envoy.HttpFilterInstance.
type bodiesHttpFilterInstance struct {
	envoyFilter envoy.EnvoyFilterInstance
}

// EventHttpRequestHeaders implements envoy.HttpFilterInstance.
func (h *bodiesHttpFilterInstance) EventHttpRequestHeaders(envoy.RequestHeaders, bool) envoy.EventHttpRequestHeadersStatus {
	return envoy.EventHttpRequestHeadersStatusContinue
}

// EventHttpRequestBody implements envoy.HttpFilterInstance.
func (h *bodiesHttpFilterInstance) EventHttpRequestBody(body envoy.RequestBodyBuffer, endOfStream bool) envoy.EventHttpRequestBodyStatus {
	fmt.Printf("new request body frame: %s\n", string(body.Copy()))
	if !endOfStream {
		// Wait for the end of the stream to see the full body.
		return envoy.EventHttpRequestBodyStatusStopIterationAndBuffer
	}

	// Now we can read the entire body.
	entireBody := h.envoyFilter.GetRequestBodyBuffer()

	// This copies the entire body into a single contiguous buffer in Go.
	fmt.Printf("entire request body: %s", string(entireBody.Copy()))

	// This demonstrates how to use ReadAt to read the body at a specific offset.
	var buf [2]byte
	for i := 0; ; i += 2 {
		_, err := entireBody.ReadAt(buf[:], int64(i))
		if err == io.EOF {
			break
		}
		fmt.Printf("request body read 2 bytes offset at %d: \"%s\"\n", i, string(buf[:]))
	}

	// Replace the entire body with 'X' without copying.
	entireBody.Slices(func(view []byte) {
		for i := 0; i < len(view); i++ {
			view[i] = 'X'
		}
	})
	return envoy.EventHttpRequestBodyStatusContinue
}

// EventHttpResponseHeaders implements envoy.HttpFilterInstance.
func (h *bodiesHttpFilterInstance) EventHttpResponseHeaders(envoy.ResponseHeaders, bool) envoy.EventHttpResponseHeadersStatus {
	return envoy.EventHttpResponseHeadersStatusContinue
}

// EventHttpResponseBody implements envoy.HttpFilterInstance.
func (h *bodiesHttpFilterInstance) EventHttpResponseBody(body envoy.ResponseBodyBuffer, endOfStream bool) envoy.EventHttpResponseBodyStatus {
	fmt.Printf("new request body frame: %s\n", string(body.Copy()))
	if !endOfStream {
		// Wait for the end of the stream to see the full body.
		return envoy.EventHttpResponseBodyStatusStopIterationAndBuffer
	}

	// Now we can read the entire body.
	entireBody := h.envoyFilter.GetResponseBodyBuffer()

	// This copies the entire body into a single contiguous buffer in Go.
	fmt.Printf("entire response body: %s", string(entireBody.Copy()))

	// This demonstrates how to use ReadAt to read the body at a specific offset.
	var buf [2]byte
	for i := 0; ; i += 2 {
		_, err := entireBody.ReadAt(buf[:], int64(i))
		if err == io.EOF {
			break
		}
		fmt.Printf("response body read 2 bytes offset at %d: \"%s\"\n", i, string(buf[:]))
	}

	// Replace the entire body with 'Y' without copying.
	entireBody.Slices(func(view []byte) {
		for i := 0; i < len(view); i++ {
			view[i] = 'Y'
		}
	})
	return envoy.EventHttpResponseBodyStatusContinue
}

// EventHttpDestroy implements envoy.HttpFilterInstance.
func (h *bodiesHttpFilterInstance) EventHttpDestroy(envoy.EnvoyFilterInstance) {}
