package main

import (
	"fmt"
	"log"

	"github.com/envoyproxyx/go-sdk/envoy"
)

// headersHttpFilter implements envoy.HttpFilter.
//
// This is to demonstrate how to use header manipulation APIs.
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
	fooValue, _ := headers.Get("foo")
	if !fooValue.Equal("value") {
		log.Fatalf("expected foo to be \"value\", got %s", fooValue.String())
	}
	fmt.Println("foo:", fooValue.String())
	headers.Values("multiple-values", func(value envoy.HeaderValue) { fmt.Println("multiple-values:", value.String()) })
	headers.Remove("multiple-values")
	headers.Set("foo", "yes")
	headers.Set("multiple-values-to-be-single", "single")
	return envoy.EventHttpRequestHeadersStatusContinue
}

// EventHttpRequestBody implements envoy.HttpFilterInstance.
func (h *headersHttpFilterInstance) EventHttpRequestBody(envoy.RequestBodyBuffer, bool) envoy.EventHttpRequestBodyStatus {
	return envoy.EventHttpRequestBodyStatusContinue
}

// EventHttpResponseHeaders implements envoy.HttpFilterInstance.
func (h *headersHttpFilterInstance) EventHttpResponseHeaders(headers envoy.ResponseHeaders, _ bool) envoy.EventHttpResponseHeadersStatus {
	headers.Values("this-is", func(value envoy.HeaderValue) {
		if !value.Equal("response-header") {
			log.Fatalf("expected this-is to be \"response-header\", got %s", value.String())
		}
		fmt.Println("this-is:", value.String())
	})
	headers.Values("this-is-2", func(value envoy.HeaderValue) { fmt.Println("this-is-2:", value.String()) })

	headers.Set("this-is", "response-header")
	headers.Remove("this-is-2")
	headers.Set("multiple-values-res-to-be-single", "single")
	return envoy.EventHttpResponseHeadersStatusContinue
}

// EventHttpResponseBody implements envoy.HttpFilterInstance.
func (h *headersHttpFilterInstance) EventHttpResponseBody(envoy.ResponseBodyBuffer, bool) envoy.EventHttpResponseBodyStatus {
	return envoy.EventHttpResponseBodyStatusContinue
}

// EventHttpDestroy implements envoy.HttpFilterInstance.
func (h *headersHttpFilterInstance) EventHttpDestroy(envoy.EnvoyFilterInstance) {}
