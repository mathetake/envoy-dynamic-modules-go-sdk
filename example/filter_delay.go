package main

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/envoyproxyx/go-sdk/envoy"
)

// delayHttpFilter implements envoy.HttpFilter.
type delayHttpFilter struct{ requestCounts atomic.Int32 }

func newDelayHttpFilter(string) envoy.HttpFilter { return &delayHttpFilter{} }

// NewHttpFilterInstance implements envoy.HttpFilter.
func (m *delayHttpFilter) NewHttpFilterInstance(e envoy.EnvoyFilterInstance) envoy.HttpFilterInstance {
	// NewHttpFilterInstance is called for each new Http request, so we can use a counter to track the number of requests.
	// On the other hand, that means this function must be thread-safe.
	id := m.requestCounts.Add(1)
	return &delayHttpFilterInstance{id: id, envoyFilter: e}
}

// Destroy implements envoy.HttpFilter.
func (m *delayHttpFilter) Destroy() {
	fmt.Println("Destroy called")
}

// delayHttpFilterInstance implements envoy.HttpFilterInstance.
type delayHttpFilterInstance struct {
	id          int32
	envoyFilter envoy.EnvoyFilterInstance
}

// EventHttpRequestHeaders implements envoy.HttpFilterInstance.
func (h *delayHttpFilterInstance) EventHttpRequestHeaders(_ envoy.RequestHeaders, _ bool) envoy.EventHttpRequestHeadersStatus {
	if h.id == 1 {
		go func() {
			fmt.Println("blocking for 1 second at EventHttpRequestHeaders with id", h.id)
			time.Sleep(1 * time.Second)
			fmt.Println("calling ContinueRequest with id", h.id)
			h.envoyFilter.ContinueRequest()
		}()
		fmt.Println("EventHttpRequestHeaders returning StopAllIterationAndBuffer with id", h.id)
		return envoy.EventHttpRequestHeadersStatusStopAllIterationAndBuffer
	}
	fmt.Println("EventHttpRequestHeaders called with id", h.id)
	return envoy.EventHttpRequestHeadersStatusContinue
}

// EventHttpRequestBody implements envoy.HttpFilterInstance.
func (h *delayHttpFilterInstance) EventHttpRequestBody(_ envoy.RequestBodyBuffer, _ bool) envoy.EventHttpRequestBodyStatus {
	if h.id == 2 {
		go func() {
			fmt.Println("blocking for 1 second at EventHttpRequestBody with id", h.id)
			time.Sleep(1 * time.Second)
			fmt.Println("calling ContinueRequest with id", h.id)
			h.envoyFilter.ContinueRequest()
		}()
		fmt.Println("EventHttpRequestBody returning StopIterationAndBuffer with id", h.id)
		return envoy.EventHttpRequestBodyStatusStopIterationAndBuffer
	}
	fmt.Println("EventHttpRequestBody called with id", h.id)
	return envoy.EventHttpRequestBodyStatusContinue
}

// EventHttpResponseHeaders implements envoy.HttpFilterInstance.
func (h *delayHttpFilterInstance) EventHttpResponseHeaders(_ envoy.ResponseHeaders, _ bool) envoy.EventHttpResponseHeadersStatus {
	if h.id == 3 {
		go func() {
			fmt.Println("blocking for 1 second at EventHttpResponseHeaders with id", h.id)
			time.Sleep(1 * time.Second)
			fmt.Println("calling ContinueResponse with id", h.id)
			h.envoyFilter.ContinueResponse()
		}()
		fmt.Println("EventHttpResponseHeaders returning StopAllIterationAndBuffer with id", h.id)
		return envoy.EventHttpResponseHeadersStatusStopAllIterationAndBuffer
	}
	fmt.Println("EventHttpResponseHeaders called with id", h.id)
	return envoy.EventHttpResponseHeadersStatusContinue
}

// EventHttpResponseBody implements envoy.HttpFilterInstance.
func (h *delayHttpFilterInstance) EventHttpResponseBody(_ envoy.ResponseBodyBuffer, _ bool) envoy.EventHttpResponseBodyStatus {
	if h.id == 4 {
		go func() {
			fmt.Println("blocking for 1 second at EventHttpResponseBody with id", h.id)
			time.Sleep(1 * time.Second)
			fmt.Println("calling ContinueResponse with id", h.id)
			h.envoyFilter.ContinueResponse()
		}()
		fmt.Println("EventHttpResponseBody returning StopIterationAndBuffer with id", h.id)
		return envoy.EventHttpResponseBodyStatusStopIterationAndBuffer
	}
	fmt.Println("EventHttpResponseBody called with id", h.id)
	return envoy.EventHttpResponseBodyStatusContinue
}

// EventHttpDestroy implements envoy.HttpFilterInstance.
func (h *delayHttpFilterInstance) EventHttpDestroy(envoy.EnvoyFilterInstance) {}
