//go:build cgo

package envoy

/*
#include "abi.h"
*/
import "C"
import (
	"unsafe"
)

// This file corresponds to the event hooks defined in https://github.com/envoyproxyx/abi/blob/main/abi.h.

// __envoy_dynamic_module_v1_event_module_init is called by the main thread when the module is
// loaded exactly once per module. The function returns 0 on success and non-zero on failure.
//
//export __envoy_dynamic_module_v1_event_module_init
func __envoy_dynamic_module_v1_event_module_init(
	configPtr C.__envoy_dynamic_module_v1_type_ModuleConfigPtr,
	configSize C.__envoy_dynamic_module_v1_type_ModuleConfigSize) C.__envoy_dynamic_module_v1_type_ModuleContextPtr {
	rawStr := unsafe.String((*byte)(unsafe.Pointer(uintptr(configPtr))), configSize)
	// Copy the config string to Go memory so that the caller can take ownership of the memory.
	var configStrCopy = make([]byte, len(rawStr))
	copy(configStrCopy, rawStr)
	// Call the exported function from the Go module.
	moduleContext := NewModuleContext(rawStr)
	memManager.pinModuleContext(moduleContext)
	return C.__envoy_dynamic_module_v1_type_ModuleContextPtr((uintptr)(unsafe.Pointer(memManager.unwrapPinnedModuleContext())))
}

// __envoy_dynamic_module_v1_event_http_context_init is called by any worker thread when a new
// stream is created. That means that the function should be thread-safe.
//
// The function returns a pointer to a new instance of the context or nullptr on failure.
// The lifetime of the returned pointer should be managed by the dynamic module.
//
//export __envoy_dynamic_module_v1_event_http_context_init
func __envoy_dynamic_module_v1_event_http_context_init(
	envoyFilterPtr C.__envoy_dynamic_module_v1_type_EnvoyFilterPtr,
	moduleCtx C.__envoy_dynamic_module_v1_type_ModuleContextPtr,
) C.__envoy_dynamic_module_v1_type_HttpContextPtr {
	envoyPtr := envoyFilterC{Raw: envoyFilterPtr}
	m := *(*ModuleContext)(unsafe.Pointer(uintptr(moduleCtx))) //nolint:govet
	httpCtx := m.HttpContextInit(envoyPtr)
	pined := memManager.pinHttpContext(httpCtx)
	return C.__envoy_dynamic_module_v1_type_HttpContextPtr(uintptr((unsafe.Pointer(pined))))
}

// __envoy_dynamic_module_v1_event_http_request_headers is called when request headers are received.
//
//export __envoy_dynamic_module_v1_event_http_request_headers
func __envoy_dynamic_module_v1_event_http_request_headers(
	envoyFilterPtr C.__envoy_dynamic_module_v1_type_EnvoyFilterPtr,
	httpContextPtr C.__envoy_dynamic_module_v1_type_HttpContextPtr,
	requestHeadersPtr C.__envoy_dynamic_module_v1_type_HttpRequestHeadersMapPtr,
	endOfStream C.__envoy_dynamic_module_v1_type_EndOfStream,
) C.__envoy_dynamic_module_v1_type_EventHttpRequestHeadersStatus {
	envoyPtr := envoyFilterC{Raw: envoyFilterPtr}
	httpCtx := unwrapRawPinHttpContext(uintptr(httpContextPtr))
	mapPtr := requestHeadersC{Raw: requestHeadersPtr}
	end := endOfStream != 0
	result := httpCtx.EventHttpRequestHeaders(envoyPtr, mapPtr, end)
	return C.__envoy_dynamic_module_v1_type_EventHttpRequestHeadersStatus(result)
}

// __envoy_dynamic_module_v1_event_http_request_body is called when request body data is received.
//
//export __envoy_dynamic_module_v1_event_http_request_body
func __envoy_dynamic_module_v1_event_http_request_body(
	envoyFilterPtr C.__envoy_dynamic_module_v1_type_EnvoyFilterPtr,
	httpContextPtr C.__envoy_dynamic_module_v1_type_HttpContextPtr,
	buffer C.__envoy_dynamic_module_v1_type_HttpRequestBodyBufferPtr,
	endOfStream C.__envoy_dynamic_module_v1_type_EndOfStream) C.__envoy_dynamic_module_v1_type_EventHttpRequestBodyStatus {
	envoyPtr := envoyFilterC{Raw: envoyFilterPtr}
	httpCtx := unwrapRawPinHttpContext(uintptr(httpContextPtr))
	buf := requestBodyBufferC{Raw: buffer}
	end := endOfStream != 0
	result := httpCtx.EventHttpRequestBody(envoyPtr, buf, end)
	return C.__envoy_dynamic_module_v1_type_EventHttpRequestBodyStatus(result)
}

// __envoy_dynamic_module_v1_event_http_response_headers is called when response headers are
// received.
//
//export __envoy_dynamic_module_v1_event_http_response_headers
func __envoy_dynamic_module_v1_event_http_response_headers(
	envoyFilterPtr C.__envoy_dynamic_module_v1_type_EnvoyFilterPtr,
	httpContextPtr C.__envoy_dynamic_module_v1_type_HttpContextPtr,
	responseHeadersMapPtr C.__envoy_dynamic_module_v1_type_HttpResponseHeaderMapPtr,
	endOfStream C.__envoy_dynamic_module_v1_type_EndOfStream) C.__envoy_dynamic_module_v1_type_EventHttpResponseHeadersStatus {
	envoyPtr := envoyFilterC{Raw: envoyFilterPtr}
	httpCtx := unwrapRawPinHttpContext(uintptr(httpContextPtr))
	mapPtr := responseHeadersC{Raw: responseHeadersMapPtr}
	end := endOfStream != 0
	result := httpCtx.EventHttpResponseHeaders(envoyPtr, mapPtr, end)
	return C.__envoy_dynamic_module_v1_type_EventHttpResponseHeadersStatus(result)
}

// __envoy_dynamic_module_v1_event_http_response_body is called when response body data is received.
//
//export __envoy_dynamic_module_v1_event_http_response_body
func __envoy_dynamic_module_v1_event_http_response_body(
	envoyFilterPtr C.__envoy_dynamic_module_v1_type_EnvoyFilterPtr,
	httpContextPtr C.__envoy_dynamic_module_v1_type_HttpContextPtr,
	buffer C.__envoy_dynamic_module_v1_type_HttpResponseBodyBufferPtr,
	endOfStream C.__envoy_dynamic_module_v1_type_EndOfStream) C.__envoy_dynamic_module_v1_type_EventHttpResponseBodyStatus {
	envoyPtr := envoyFilterC{Raw: envoyFilterPtr}
	httpCtx := unwrapRawPinHttpContext(uintptr(httpContextPtr))
	buf := responseBodyBufferC{Raw: buffer}
	end := endOfStream != 0
	result := httpCtx.EventHttpResponseBody(envoyPtr, buf, end)
	return C.__envoy_dynamic_module_v1_type_EventHttpResponseBodyStatus(result)
}

// __envoy_dynamic_module_v1_event_http_destroy is called when the stream is destroyed.
//
//export __envoy_dynamic_module_v1_event_http_destroy
func __envoy_dynamic_module_v1_event_http_destroy(
	envoyFilterPtr C.__envoy_dynamic_module_v1_type_EnvoyFilterPtr,
	httpContextPtr C.__envoy_dynamic_module_v1_type_HttpContextPtr) {
	envoyPtr := envoyFilterC{Raw: envoyFilterPtr}
	httpCtx := unwrapRawPinHttpContext(uintptr(httpContextPtr))
	httpCtx.EventHttpDestroy(envoyPtr)
	memManager.removeHttpContext((*pinedHttpContext)(unsafe.Pointer(uintptr(httpContextPtr))))
}

type envoyFilterC struct {
	Raw C.__envoy_dynamic_module_v1_type_EnvoyFilterPtr
}

func (c envoyFilterC) ContinueRequest() {
	C.__envoy_dynamic_module_v1_http_continue_request(c.Raw)
}

func (c envoyFilterC) ContinueResponse() {
	C.__envoy_dynamic_module_v1_http_continue_response(c.Raw)
}

type requestHeadersC struct {
	Raw C.__envoy_dynamic_module_v1_type_HttpRequestHeadersMapPtr
}

type responseHeadersC struct {
	Raw C.__envoy_dynamic_module_v1_type_HttpResponseHeaderMapPtr
}

type requestBodyBufferC struct {
	Raw C.__envoy_dynamic_module_v1_type_HttpRequestBodyBufferPtr
}

type responseBodyBufferC struct {
	Raw C.__envoy_dynamic_module_v1_type_HttpResponseBodyBufferPtr
}
