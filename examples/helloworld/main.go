package main

import (
	"fmt"

	"github.com/envoyproxyx/go-sdk/envoy"
)

func main() {} // main function must be present but empty.

func init() {
	// Set the envoy.NewModuleContext function to create a new module context.
	envoy.NewModuleContext = newModuleContext
}

// moduleContext implements envoy.ModuleContext.
type moduleContext struct{}

func newModuleContext(config string) envoy.ModuleContext {
	fmt.Println("NewModuleContext called:", config)
	return &moduleContext{}
}

// HttpContextInit implements envoy.ModuleContext and is called for each new Http request.
func (m *moduleContext) HttpContextInit(envoy.EnvoyFilter) envoy.HttpContext {
	fmt.Println("HttpContextInit called")
	return &httpContext{}
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
func (h httpContext) EventHttpDestroy(envoy.EnvoyFilter) {
	fmt.Println("EventHttpDestroy called")
}
