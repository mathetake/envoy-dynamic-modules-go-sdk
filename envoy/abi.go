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
	envoyPtr := &envoyFilterC{raw: envoyFilterPtr}
	m := *(*ModuleContext)(unsafe.Pointer(uintptr(moduleCtx))) //nolint:govet
	httpCtx := m.HttpContextInit(envoyPtr)
	pined := memManager.pinHttpContext(httpCtx)
	pined.envoyFilter = envoyPtr
	return C.__envoy_dynamic_module_v1_type_HttpContextPtr(uintptr((unsafe.Pointer(pined))))
}

// __envoy_dynamic_module_v1_event_http_request_headers is called when request headers are received.
//
//export __envoy_dynamic_module_v1_event_http_request_headers
func __envoy_dynamic_module_v1_event_http_request_headers(
	httpContextPtr C.__envoy_dynamic_module_v1_type_HttpContextPtr,
	requestHeadersPtr C.__envoy_dynamic_module_v1_type_HttpRequestHeadersMapPtr,
	endOfStream C.__envoy_dynamic_module_v1_type_EndOfStream,
) C.__envoy_dynamic_module_v1_type_EventHttpRequestHeadersStatus {
	httpCtx := unwrapRawPinHttpContext(uintptr(httpContextPtr))
	mapPtr := requestHeadersC{Raw: requestHeadersPtr}
	end := endOfStream != 0
	result := httpCtx.ctx.EventHttpRequestHeaders(mapPtr, end)
	return C.__envoy_dynamic_module_v1_type_EventHttpRequestHeadersStatus(result)
}

// __envoy_dynamic_module_v1_event_http_request_body is called when request body data is received.
//
//export __envoy_dynamic_module_v1_event_http_request_body
func __envoy_dynamic_module_v1_event_http_request_body(
	httpContextPtr C.__envoy_dynamic_module_v1_type_HttpContextPtr,
	buffer C.__envoy_dynamic_module_v1_type_HttpRequestBodyBufferPtr,
	endOfStream C.__envoy_dynamic_module_v1_type_EndOfStream) C.__envoy_dynamic_module_v1_type_EventHttpRequestBodyStatus {
	httpCtx := unwrapRawPinHttpContext(uintptr(httpContextPtr))
	buf := requestBodyBufferC{Raw: buffer}
	end := endOfStream != 0
	result := httpCtx.ctx.EventHttpRequestBody(buf, end)
	return C.__envoy_dynamic_module_v1_type_EventHttpRequestBodyStatus(result)
}

// __envoy_dynamic_module_v1_event_http_response_headers is called when response headers are
// received.
//
//export __envoy_dynamic_module_v1_event_http_response_headers
func __envoy_dynamic_module_v1_event_http_response_headers(
	httpContextPtr C.__envoy_dynamic_module_v1_type_HttpContextPtr,
	responseHeadersMapPtr C.__envoy_dynamic_module_v1_type_HttpResponseHeaderMapPtr,
	endOfStream C.__envoy_dynamic_module_v1_type_EndOfStream) C.__envoy_dynamic_module_v1_type_EventHttpResponseHeadersStatus {
	httpCtx := unwrapRawPinHttpContext(uintptr(httpContextPtr))
	mapPtr := responseHeadersC{Raw: responseHeadersMapPtr}
	end := endOfStream != 0
	result := httpCtx.ctx.EventHttpResponseHeaders(mapPtr, end)
	return C.__envoy_dynamic_module_v1_type_EventHttpResponseHeadersStatus(result)
}

// __envoy_dynamic_module_v1_event_http_response_body is called when response body data is received.
//
//export __envoy_dynamic_module_v1_event_http_response_body
func __envoy_dynamic_module_v1_event_http_response_body(
	httpContextPtr C.__envoy_dynamic_module_v1_type_HttpContextPtr,
	buffer C.__envoy_dynamic_module_v1_type_HttpResponseBodyBufferPtr,
	endOfStream C.__envoy_dynamic_module_v1_type_EndOfStream) C.__envoy_dynamic_module_v1_type_EventHttpResponseBodyStatus {
	httpCtx := unwrapRawPinHttpContext(uintptr(httpContextPtr))
	buf := responseBodyBufferC{Raw: buffer}
	end := endOfStream != 0
	result := httpCtx.ctx.EventHttpResponseBody(buf, end)
	return C.__envoy_dynamic_module_v1_type_EventHttpResponseBodyStatus(result)
}

// __envoy_dynamic_module_v1_event_http_destroy is called when the stream is destroyed.
//
//export __envoy_dynamic_module_v1_event_http_destroy
func __envoy_dynamic_module_v1_event_http_destroy(
	httpContextPtr C.__envoy_dynamic_module_v1_type_HttpContextPtr) {
	httpCtx := unwrapRawPinHttpContext(uintptr(httpContextPtr))
	httpCtx.ctx.EventHttpDestroy(httpCtx.envoyFilter)
	httpCtx.envoyFilter.(*envoyFilterC).destroyed = true
	memManager.removeHttpContext((*pinedHttpContext)(unsafe.Pointer(uintptr(httpContextPtr))))
}

type envoyFilterC struct {
	raw       C.__envoy_dynamic_module_v1_type_EnvoyFilterPtr
	destroyed bool
}

func (c *envoyFilterC) ContinueRequest() {
	if c.destroyed {
		return
	}
	C.__envoy_dynamic_module_v1_http_continue_request(c.raw)
}

func (c *envoyFilterC) ContinueResponse() {
	if c.destroyed {
		return
	}
	C.__envoy_dynamic_module_v1_http_continue_response(c.raw)
}

func (c *envoyFilterC) Destroyed() bool {
	return c.destroyed
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
