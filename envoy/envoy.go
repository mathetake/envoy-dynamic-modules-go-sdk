package envoy

// NewModuleContext is a function that creates a new ModuleContext.
// This is a global variable that should be set in the main function.
//
// The function is called once globally. The function is only called by the main thread,
// so it does not need to be thread-safe.
//
// `config` is the configuration string that is passed to the module that is set in the Envoy configuration.
var NewModuleContext func(config string) ModuleContext

// EnvoyFilterPtr is an opaque object that represents the underlying Envoy Http filter instance.
// This is used to interact with it from the module code.
type EnvoyFilterPtr struct {
	Raw __envoy_dynamic_module_v1_type_EnvoyFilterPtr
}

// RequestHeadersMapPtr is an opaque object that represents the underlying Envoy Http request headers map.
// This is used to interact with it from the module code.
type RequestHeadersMapPtr struct {
	Raw __envoy_dynamic_module_v1_type_HttpRequestHeadersMapPtr
}

// ResponseHeadersMapPtr is an opaque object that represents the underlying Envoy Http response headers map.
// This is used to interact with it from the module code.
type ResponseHeadersMapPtr struct {
	Raw __envoy_dynamic_module_v1_type_HttpResponseHeaderMapPtr
}

// RequestBodyBufferPtr is an opaque object that represents the underlying Envoy Http request body buffer.
// This is used to interact with it from the module code.
type RequestBodyBufferPtr struct {
	Raw __envoy_dynamic_module_v1_type_HttpRequestBodyBufferPtr
}

// ResponseBodyBufferPtr is an opaque object that represents the underlying Envoy Http response body buffer.
// This is used to interact with it from the module code.
type ResponseBodyBufferPtr struct {
	Raw __envoy_dynamic_module_v1_type_HttpResponseBodyBufferPtr
}

// ModuleContext is an interface that represents the module context.
// It is used to create HttpContext objects that correspond to each Http request.
//
// This is only created once per module instance via the NewModuleContext function.
type ModuleContext interface {
	// HttpContextInit is called for each new Http request.
	// Note that this must be concurrency-safe as it can be called concurrently for multiple requests.
	HttpContextInit(EnvoyFilterPtr) HttpContext
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
	EventHttpRequestHeaders(EnvoyFilterPtr, RequestHeadersMapPtr, bool) EventHttpRequestHeadersStatus
	// EventHttpRequestBody is called when request body data is received.
	// The function should return the status of the operation.
	//
	//  * `envoyFilter` is the pointer to the Envoy filter instance.
	//  * `requestBody` is the pointer to the request body buffer.
	//  * `endOfStream` is a boolean that indicates if this is the last data frame.
	EventHttpRequestBody(EnvoyFilterPtr, RequestBodyBufferPtr, bool) EventHttpRequestBodyStatus
	// EventHttpResponseHeaders is called when response headers are received.
	// The function should return the status of the operation.
	//
	//  * `envoyFilter` is the pointer to the Envoy filter instance.
	//  * `responseHeaders` is the pointer to the response headers map.
	//  * `endOfStream` is a boolean that indicates if this is the headers-only response.
	EventHttpResponseHeaders(EnvoyFilterPtr, ResponseHeadersMapPtr, bool) EventHttpResponseHeadersStatus
	// EventHttpResponseBody is called when response body data is received.
	// The function should return the status of the operation.
	//
	//  * `envoyFilter` is the pointer to the Envoy filter instance.
	//  * `responseBody` is the pointer to the response body buffer.
	//  * `endOfStream` is a boolean that indicates if this is the last data frame.
	EventHttpResponseBody(EnvoyFilterPtr, ResponseBodyBufferPtr, bool) EventHttpResponseBodyStatus

	// EventHttpDestroy is called when the stream is destroyed.
	// This is called when the stream is completed or when the stream is reset.
	EventHttpDestroy(EnvoyFilterPtr)
}

// EventHttpRequestHeadersStatus is the return value of the HttpContext.EventHttpRequestHeaders event.
type EventHttpRequestHeadersStatus int

const (
	// EventHttpRequestHeadersStatusContinue is returned when the operation should continue.
	EventHttpRequestHeadersStatusContinue EventHttpRequestHeadersStatus = iota
)

// EventHttpRequestBodyStatus is the return value of the HttpContext.EventHttpRequestBody event.
type EventHttpRequestBodyStatus int

const (
	// EventHttpRequestBodyStatusContinue is returned when the operation should continue.
	EventHttpRequestBodyStatusContinue EventHttpRequestBodyStatus = iota
)

// EventHttpResponseHeadersStatus is the return value of the HttpContext.EventHttpResponseHeaders event.
type EventHttpResponseHeadersStatus int

const (
	// EventHttpResponseHeadersStatusContinue is returned when the operation should continue.
	EventHttpResponseHeadersStatusContinue EventHttpResponseHeadersStatus = iota
)

// EventHttpResponseBodyStatus is the return value of the HttpContext.EventHttpResponseBody event.
type EventHttpResponseBodyStatus int

const (
	// EventHttpResponseBodyStatusContinue is returned when the operation should continue.
	EventHttpResponseBodyStatusContinue EventHttpResponseBodyStatus = iota
)
