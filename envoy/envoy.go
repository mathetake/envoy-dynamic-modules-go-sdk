package envoy

// NewModuleContext is a function that creates a new ModuleContext.
// This is a global variable that should be set in the main function.
//
// The function is called once globally. The function is only called by the main thread,
// so it does not need to be thread-safe.
//
// `config` is the configuration string that is passed to the module that is set in the Envoy configuration.
var NewModuleContext func(config string) ModuleContext

// EnvoyFilter is an opaque object that represents the underlying Envoy Http filter instance.
// This is used to interact with it from the module code.
type EnvoyFilter interface {
	// ContinueRequest is a function that continues the request processing.
	ContinueRequest()
	// ContinueResponse is a function that continues the response processing.
	ContinueResponse()
}

// RequestHeaders is an opaque object that represents the underlying Envoy Http request headers map.
// This is used to interact with it from the module code.
type RequestHeaders interface{}

// ResponseHeadersMap is an opaque object that represents the underlying Envoy Http response headers map.
// This is used to interact with it from the module code.
type ResponseHeaders interface{}

// RequestBodyBuffer is an opaque object that represents the underlying Envoy Http request body buffer.
// This is used to interact with it from the module code.
type RequestBodyBuffer interface{}

// ResponseBodyBuffer is an opaque object that represents the underlying Envoy Http response body buffer.
// This is used to interact with it from the module code.
type ResponseBodyBuffer interface{}

// ModuleContext is an interface that represents the module context.
// It is used to create HttpContext objects that correspond to each Http request.
//
// This is only created once per module instance via the NewModuleContext function.
type ModuleContext interface {
	// HttpContextInit is called for each new Http request.
	// Note that this must be concurrency-safe as it can be called concurrently for multiple requests.
	HttpContextInit(EnvoyFilter) HttpContext
}

// HttpContext is an interface that represents the Http request context.
//
// The context is created for each new Http request and is destroyed when the request is completed.
type HttpContext interface {
	// EventHttpRequestHeaders is called when request headers are received.
	// The function should return the status of the operation.
	//
	//  * `envoyFilter` is the pointer to the Envoy filter instance.
	//  * `requestHeaders` is the pointer to the request headers map.
	//  * `endOfStream` is a boolean that indicates if this is the headers-only request.
	EventHttpRequestHeaders(EnvoyFilter, RequestHeaders, bool) EventHttpRequestHeadersStatus
	// EventHttpRequestBody is called when request body data is received.
	// The function should return the status of the operation.
	//
	//  * `envoyFilter` is the pointer to the Envoy filter instance.
	//  * `requestBody` is the pointer to the request body buffer.
	//  * `endOfStream` is a boolean that indicates if this is the last data frame.
	EventHttpRequestBody(EnvoyFilter, RequestBodyBuffer, bool) EventHttpRequestBodyStatus
	// EventHttpResponseHeaders is called when response headers are received.
	// The function should return the status of the operation.
	//
	//  * `envoyFilter` is the pointer to the Envoy filter instance.
	//  * `responseHeaders` is the pointer to the response headers map.
	//  * `endOfStream` is a boolean that indicates if this is the headers-only response.
	EventHttpResponseHeaders(EnvoyFilter, ResponseHeaders, bool) EventHttpResponseHeadersStatus
	// EventHttpResponseBody is called when response body data is received.
	// The function should return the status of the operation.
	//
	//  * `envoyFilter` is the pointer to the Envoy filter instance.
	//  * `responseBody` is the pointer to the response body buffer.
	//  * `endOfStream` is a boolean that indicates if this is the last data frame.
	EventHttpResponseBody(EnvoyFilter, ResponseBodyBuffer, bool) EventHttpResponseBodyStatus

	// EventHttpDestroy is called when the stream is destroyed.
	// This is called when the stream is completed or when the stream is reset.
	EventHttpDestroy(EnvoyFilter)
}
