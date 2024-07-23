package envoy

import "unsafe"

// NewHttpFilter is a function that creates a new HttpFilter that corresponds to a filter in the Envoy filter chain.
// This is a global variable that should be set in the init function in the program once.
//
// The function is called once globally. The function is only called by the main thread,
// so it does not need to be thread-safe.
//
// `config` is the configuration string that is passed to the module that is set in the Envoy configuration.
var NewHttpFilter func(config string) HttpFilter

// EnvoyFilterInstance is an opaque object that represents the underlying Envoy Http filter instance.
// This is used to interact with it from the module code.
type EnvoyFilterInstance interface {
	// ContinueRequest is a function that continues the request processing.
	ContinueRequest()
	// ContinueResponse is a function that continues the response processing.
	ContinueResponse()
	// Destroyed returns true if the stream is destroyed.
	Destroyed() bool
}

// RequestHeaders is an opaque object that represents the underlying Envoy Http request headers map.
// This is used to interact with it from the module code.
type RequestHeaders interface {
	// Get iterates over the header values for the given key.
	// Usually, there is only one value for a key, but in general, multiple values are possible.
	Get(key string, iter func(value HeaderValue))
	// Set sets the value for the given key. If multiple values are set for the same key,
	// this removes all the previous values and sets the new single value.
	Set(key, value string)
	// Remove removes the value for the given key. If multiple values are set for the same key,
	// this removes all the values.
	Remove(key string)
}

// ResponseHeadersMap is an opaque object that represents the underlying Envoy Http response headers map.
// This is used to interact with it from the module code.
type ResponseHeaders interface {
	// Get iterates over the header values for the given key.
	// Usually, there is only one value for a key, but in general, multiple values are possible.
	Get(key string, iter func(value HeaderValue))
	// Set sets the value for the given key. If multiple values are set for the same key,
	// this removes all the previous values and sets the new single value.
	Set(key, value string)
	// Remove removes the value for the given key. If multiple values are set for the same key,
	// this removes all the values.
	Remove(key string)
}

// HeaderValue represents a single header value whose data is owned by the Envoy.
type HeaderValue struct {
	data *byte
	size int
}

// String returns the string representation of the header value.
// This copies the underlying data to a new buffer and returns the string.
func (h HeaderValue) String() string {
	view := unsafe.Slice(h.data, h.size)
	return string(view)
}

// Equal returns true if the header value is equal to the given string.
func (h HeaderValue) Equal(str string) bool {
	v := unsafe.String(h.data, h.size)
	return v == str
}

// RequestBodyBuffer is an opaque object that represents the underlying Envoy Http request body buffer.
// This is used to interact with it from the module code.
type RequestBodyBuffer interface{}

// ResponseBodyBuffer is an opaque object that represents the underlying Envoy Http response body buffer.
// This is used to interact with it from the module code.
type ResponseBodyBuffer interface{}

// HttpFilter is an interface that represents a single http filter in the Envoy filter chain.
// It is used to create HttpFilterInstance(s) that correspond to each Http request.
//
// This is only created once per module instance via the NewHttpFilter function.
type HttpFilter interface {
	// NewHttpFilterInstance is called for each new Http request.
	// Note that this must be concurrency-safe as it can be called concurrently for multiple requests.
	//
	// * `EnvoyFilterInstance` is the Envoy filter object that is used to interact with the underlying Envoy filter.
	//  This object is unique for each Http request. The object is destroyed when the stream is destroyed.
	//  Therefore, after EventHttpDestroy is called, the methods on this object become no-op.
	NewHttpFilterInstance(EnvoyFilterInstance) HttpFilterInstance

	// Destroy is called when this filter is destroyed. E.g. the filter chain configuration is updated and removed from the Envoy.
	Destroy()
}

// HttpFilterInstance is an interface that represents each Http request.
//
// Thisis created for each new Http request and is destroyed when the request is completed.
type HttpFilterInstance interface {
	// EventHttpRequestHeaders is called when request headers are received.
	// The function should return the status of the operation.
	//
	//  * `requestHeaders` is the pointer to the request headers map.
	//  * `endOfStream` is a boolean that indicates if this is the headers-only request.
	EventHttpRequestHeaders(RequestHeaders, bool) EventHttpRequestHeadersStatus
	// EventHttpRequestBody is called when request body data is received.
	// The function should return the status of the operation.
	//
	//  * `requestBody` is the pointer to the request body buffer.
	//  * `endOfStream` is a boolean that indicates if this is the last data frame.
	EventHttpRequestBody(RequestBodyBuffer, bool) EventHttpRequestBodyStatus
	// EventHttpResponseHeaders is called when response headers are received.
	// The function should return the status of the operation.
	//
	//  * `responseHeaders` is the pointer to the response headers map.
	//  * `endOfStream` is a boolean that indicates if this is the headers-only response.
	EventHttpResponseHeaders(ResponseHeaders, bool) EventHttpResponseHeadersStatus
	// EventHttpResponseBody is called when response body data is received.
	// The function should return the status of the operation.
	//
	//  * `responseBody` is the pointer to the response body buffer.
	//  * `endOfStream` is a boolean that indicates if this is the last data frame.
	EventHttpResponseBody(ResponseBodyBuffer, bool) EventHttpResponseBodyStatus

	// EventHttpDestroy is called when the stream is destroyed.
	// This is called when the stream is completed or when the stream is reset.
	EventHttpDestroy(EnvoyFilterInstance)
}
