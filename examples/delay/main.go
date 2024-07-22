package main

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/envoyproxyx/go-sdk/envoy"
)

func main() {} // main function must be present but empty.

func init() {
	// Set the envoy.NewModuleContext function to create a new module context.
	envoy.NewModuleContext = newModuleContext
}

// moduleContext implements envoy.ModuleContext.
type moduleContext struct{ requestCounts atomic.Int32 }

func newModuleContext(config string) envoy.ModuleContext { return &moduleContext{} }

// HttpContextInit implements envoy.ModuleContext and is called for each new Http request.
func (m *moduleContext) HttpContextInit(e envoy.EnvoyFilter) envoy.HttpContext {
	// HttpContextInit is called for each new Http request, so we can use a counter to track the number of requests.
	// On the other hand, that means this function must be thread-safe.
	id := m.requestCounts.Add(1)
	return &httpContext{id: id, envoyFilter: e}
}

// httpContext implements envoy.HttpContext.
type httpContext struct {
	id          int32
	envoyFilter envoy.EnvoyFilter
}

// EventHttpRequestHeaders implements envoy.HttpContext.
func (h httpContext) EventHttpRequestHeaders(_ envoy.RequestHeaders, _ bool) envoy.EventHttpRequestHeadersStatus {
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

// EventHttpRequestBody implements envoy.HttpContext.
func (h *httpContext) EventHttpRequestBody(_ envoy.RequestBodyBuffer, _ bool) envoy.EventHttpRequestBodyStatus {
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

// EventHttpResponseHeaders implements envoy.HttpContext.
func (h *httpContext) EventHttpResponseHeaders(_ envoy.ResponseHeaders, _ bool) envoy.EventHttpResponseHeadersStatus {
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

// EventHttpResponseBody implements envoy.HttpContext.
func (h *httpContext) EventHttpResponseBody(_ envoy.ResponseBodyBuffer, _ bool) envoy.EventHttpResponseBodyStatus {
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

// EventHttpDestroy implements envoy.HttpContext.
func (h *httpContext) EventHttpDestroy(envoy.EnvoyFilter) {}
