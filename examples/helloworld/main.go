package main

import (
	"fmt"

	"github.com/envoyproxyx/go-sdk/envoy"
)

func main() {} // main function must be present but empty.

func init() {
	// Set the envoy.NewHttpFilter function to create a new module context.
	envoy.NewHttpFilter = newHttpFilter
}

// httpFilter implements envoy.HttpFilter.
type httpFilter struct{}

func newHttpFilter(config string) envoy.HttpFilter {
	fmt.Println("NewHttpFilter called:", config)
	return &httpFilter{}
}

// NewHttpFilterInstance implements envoy.HttpFilter.
func (m *httpFilter) NewHttpFilterInstance(envoy.EnvoyFilterInstance) envoy.HttpFilterInstance {
	fmt.Println("NewHttpFilterInstance called")
	return &httpContext{}
}

// Destroy implements envoy.HttpContext.
func (m *httpFilter) Destroy() {
	fmt.Println("Destroy called")
}

// httpContext implements envoy.HttpContext.
type httpContext struct{}

// EventHttpRequestHeaders implements envoy.HttpContext.
func (h httpContext) EventHttpRequestHeaders(envoy.RequestHeaders, bool) envoy.EventHttpRequestHeadersStatus {
	fmt.Println("EventHttpRequestHeaders called")
	return envoy.EventHttpRequestHeadersStatusContinue
}

// EventHttpRequestBody implements envoy.HttpContext.
func (h httpContext) EventHttpRequestBody(envoy.RequestBodyBuffer, bool) envoy.EventHttpRequestBodyStatus {
	fmt.Println("EventHttpRequestBody called")
	return envoy.EventHttpRequestBodyStatusContinue
}

// EventHttpResponseHeaders implements envoy.HttpContext.
func (h httpContext) EventHttpResponseHeaders(envoy.ResponseHeaders, bool) envoy.EventHttpResponseHeadersStatus {
	fmt.Println("EventHttpResponseHeaders called")
	return envoy.EventHttpResponseHeadersStatusContinue
}

// EventHttpResponseBody implements envoy.HttpContext.
func (h httpContext) EventHttpResponseBody(envoy.ResponseBodyBuffer, bool) envoy.EventHttpResponseBodyStatus {
	fmt.Println("EventHttpResponseBody called")
	return envoy.EventHttpResponseBodyStatusContinue
}

// EventHttpDestroy implements envoy.HttpContext.
func (h httpContext) EventHttpDestroy(envoy.EnvoyFilterInstance) {
	fmt.Println("EventHttpDestroy called")
}
