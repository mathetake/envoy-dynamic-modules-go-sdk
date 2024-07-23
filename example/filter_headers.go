package main

import (
	"fmt"

	"github.com/envoyproxyx/go-sdk/envoy"
)

type headersHttpFilter struct{}

func newHeadersHttpFilter(string) envoy.HttpFilter { return &headersHttpFilter{} }

// NewHttpFilterInstance implements envoy.HttpFilter.
func (f *headersHttpFilter) NewHttpFilterInstance(envoy.EnvoyFilterInstance) envoy.HttpFilterInstance {
	return &headersHttpFilterInstance{}
}

// Destroy implements envoy.HttpFilter.
func (f *headersHttpFilter) Destroy() {}

// headersHttpFilterInstance implements envoy.HttpFilterInstance.
type headersHttpFilterInstance struct{}

// EventHttpRequestHeaders implements envoy.HttpFilterInstance.
func (h *headersHttpFilterInstance) EventHttpRequestHeaders(headers envoy.RequestHeaders, _ bool) envoy.EventHttpRequestHeadersStatus {
	headers.Get("foo", func(value envoy.HeaderValue) { fmt.Println("foo:", value.String()) })
	headers.Get("multiple-values", func(value envoy.HeaderValue) { fmt.Println("multiple-values:", value.String()) })
	return envoy.EventHttpRequestHeadersStatusContinue
}

// EventHttpRequestBody implements envoy.HttpFilterInstance.
func (h *headersHttpFilterInstance) EventHttpRequestBody(envoy.RequestBodyBuffer, bool) envoy.EventHttpRequestBodyStatus {
	return envoy.EventHttpRequestBodyStatusContinue
}

// EventHttpResponseHeaders implements envoy.HttpFilterInstance.
func (h *headersHttpFilterInstance) EventHttpResponseHeaders(headers envoy.ResponseHeaders, _ bool) envoy.EventHttpResponseHeadersStatus {
	headers.Get("this-is", func(value envoy.HeaderValue) { fmt.Println("this-is:", value.String()) })
	headers.Get("this-is-2", func(value envoy.HeaderValue) { fmt.Println("this-is-2:", value.String()) })

	return envoy.EventHttpResponseHeadersStatusContinue
}

// EventHttpResponseBody implements envoy.HttpFilterInstance.
func (h *headersHttpFilterInstance) EventHttpResponseBody(envoy.ResponseBodyBuffer, bool) envoy.EventHttpResponseBodyStatus {
	return envoy.EventHttpResponseBodyStatusContinue
}

// EventHttpDestroy implements envoy.HttpFilterInstance.
func (h *headersHttpFilterInstance) EventHttpDestroy(envoy.EnvoyFilterInstance) {}
