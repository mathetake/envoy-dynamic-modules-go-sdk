package main

import (
	"fmt"

	"github.com/envoyproxyx/go-sdk/envoy"
)

type helloWorldFilter struct{}

func newHelloWorldHttpFilter(string) envoy.HttpFilter {
	fmt.Println("newHelloWorldHttpFilter called")
	return &helloWorldFilter{}
}

// NewHttpFilterInstance implements envoy.HttpFilter.
func (f *helloWorldFilter) NewHttpFilterInstance(envoy.EnvoyFilterInstance) envoy.HttpFilterInstance {
	fmt.Println("helloWorldHttpFilter.NewHttpFilterInstance called")
	return &helloWorldHttpFilterInstance{}
}

// Destroy implements envoy.HttpContext.
func (f *helloWorldFilter) Destroy() { fmt.Println("helloWorldFilter.Destroy called") }

// helloWorldHttpFilterInstance implements envoy.HttpFilterInstance.
type helloWorldHttpFilterInstance struct{}

// EventHttpRequestHeaders implements envoy.HttpFilterInstance.
func (h *helloWorldHttpFilterInstance) EventHttpRequestHeaders(envoy.RequestHeaders, bool) envoy.EventHttpRequestHeadersStatus {
	fmt.Println("helloWorldHttpFilterInstance.EventHttpRequestHeaders called")
	return envoy.EventHttpRequestHeadersStatusContinue
}

// EventHttpRequestBody implements envoy.HttpFilterInstance.
func (h *helloWorldHttpFilterInstance) EventHttpRequestBody(envoy.RequestBodyBuffer, bool) envoy.EventHttpRequestBodyStatus {
	fmt.Println("helloWorldHttpFilterInstance.EventHttpRequestBody called")
	return envoy.EventHttpRequestBodyStatusContinue
}

// EventHttpResponseHeaders implements envoy.HttpFilterInstance.
func (h *helloWorldHttpFilterInstance) EventHttpResponseHeaders(envoy.ResponseHeaders, bool) envoy.EventHttpResponseHeadersStatus {
	fmt.Println("helloWorldHttpFilterInstance.EventHttpResponseHeaders called")
	return envoy.EventHttpResponseHeadersStatusContinue
}

// EventHttpResponseBody implements envoy.HttpFilterInstance.
func (h *helloWorldHttpFilterInstance) EventHttpResponseBody(envoy.ResponseBodyBuffer, bool) envoy.EventHttpResponseBodyStatus {
	fmt.Println("helloWorldHttpFilterInstance.EventHttpResponseBody called")
	return envoy.EventHttpResponseBodyStatusContinue
}

// EventHttpDestroy implements envoy.HttpFilterInstance.
func (h *helloWorldHttpFilterInstance) EventHttpDestroy(envoy.EnvoyFilterInstance) {
	fmt.Println("helloWorldHttpFilterInstance.EventHttpDestroy called")
}
