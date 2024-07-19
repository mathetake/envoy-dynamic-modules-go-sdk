package envoy

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
	configPtr __envoy_dynamic_module_v1_type_ModuleConfigPtr,
	configSize __envoy_dynamic_module_v1_type_ModuleConfigSize) __envoy_dynamic_module_v1_type_ModuleContextPtr {
	rawStr := unsafe.String(configPtr, configSize)
	// Copy the config string to Go memory so that the caller can take ownership of the memory.
	var configStrCopy = make([]byte, len(rawStr))
	copy(configStrCopy, rawStr)
	// Call the exported function from the Go module.
	moduleContext := NewModuleContext(rawStr)
	memManager.pinModuleContext(moduleContext)
	return __envoy_dynamic_module_v1_type_ModuleContextPtr(unsafe.Pointer(memManager.unwrapPinnedModuleContext()))
}

// __envoy_dynamic_module_v1_event_http_context_init is called by any worker thread when a new
// stream is created. That means that the function should be thread-safe.
//
// The function returns a pointer to a new instance of the context or nullptr on failure.
// The lifetime of the returned pointer should be managed by the dynamic module.
//
//export __envoy_dynamic_module_v1_event_http_context_init
func __envoy_dynamic_module_v1_event_http_context_init(
	envoyFilterPtr __envoy_dynamic_module_v1_type_EnvoyFilterPtr,
	moduleCtx __envoy_dynamic_module_v1_type_ModuleContextPtr,
) __envoy_dynamic_module_v1_type_HttpContextPtr {
	envoyPtr := EnvoyFilterPtr{Raw: envoyFilterPtr}
	m := *(*ModuleContext)(unsafe.Pointer(moduleCtx)) //nolint:govet
	httpCtx := m.HttpContextInit(envoyPtr)
	pined := memManager.pinHttpContext(httpCtx)
	return __envoy_dynamic_module_v1_type_HttpContextPtr(unsafe.Pointer(pined))
}

// __envoy_dynamic_module_v1_event_http_request_headers is called when request headers are received.
//
//export __envoy_dynamic_module_v1_event_http_request_headers
func __envoy_dynamic_module_v1_event_http_request_headers(
	envoyFilterPtr __envoy_dynamic_module_v1_type_EnvoyFilterPtr,
	httpContextPtr __envoy_dynamic_module_v1_type_HttpContextPtr,
	requestHeadersPtr __envoy_dynamic_module_v1_type_HttpRequestHeadersMapPtr,
	endOfStream __envoy_dynamic_module_v1_type_EndOfStream,
) __envoy_dynamic_module_v1_type_EventHttpRequestHeadersStatus {
	envoyPtr := EnvoyFilterPtr{Raw: envoyFilterPtr}
	httpCtx := unwrapRawPinHttpContext(uintptr(httpContextPtr))
	mapPtr := RequestHeadersMapPtr{Raw: requestHeadersPtr}
	end := endOfStream != 0
	result := httpCtx.EventHttpRequestHeaders(envoyPtr, mapPtr, end)
	return __envoy_dynamic_module_v1_type_EventHttpRequestHeadersStatus(result)
}

// __envoy_dynamic_module_v1_event_http_request_body is called when request body data is received.
//
//export __envoy_dynamic_module_v1_event_http_request_body
func __envoy_dynamic_module_v1_event_http_request_body(
	envoyFilterPtr __envoy_dynamic_module_v1_type_EnvoyFilterPtr,
	httpContextPtr __envoy_dynamic_module_v1_type_HttpContextPtr,
	buffer __envoy_dynamic_module_v1_type_HttpRequestBodyBufferPtr,
	endOfStream __envoy_dynamic_module_v1_type_EndOfStream) __envoy_dynamic_module_v1_type_EventHttpRequestBodyStatus {
	envoyPtr := EnvoyFilterPtr{Raw: envoyFilterPtr}
	httpCtx := unwrapRawPinHttpContext(uintptr(httpContextPtr))
	buf := RequestBodyBufferPtr{Raw: buffer}
	end := endOfStream != 0
	result := httpCtx.EventHttpRequestBody(envoyPtr, buf, end)
	return __envoy_dynamic_module_v1_type_EventHttpRequestBodyStatus(result)
}

// __envoy_dynamic_module_v1_event_http_response_headers is called when response headers are
// received.
//
//export __envoy_dynamic_module_v1_event_http_response_headers
func __envoy_dynamic_module_v1_event_http_response_headers(
	envoyFilterPtr __envoy_dynamic_module_v1_type_EnvoyFilterPtr,
	httpContextPtr __envoy_dynamic_module_v1_type_HttpContextPtr,
	responseHeadersMapPtr __envoy_dynamic_module_v1_type_HttpResponseHeaderMapPtr,
	endOfStream __envoy_dynamic_module_v1_type_EndOfStream) __envoy_dynamic_module_v1_type_EventHttpResponseHeadersStatus {
	envoyPtr := EnvoyFilterPtr{Raw: envoyFilterPtr}
	httpCtx := unwrapRawPinHttpContext(uintptr(httpContextPtr))
	mapPtr := ResponseHeadersMapPtr{Raw: responseHeadersMapPtr}
	end := endOfStream != 0
	result := httpCtx.EventHttpResponseHeaders(envoyPtr, mapPtr, end)
	return __envoy_dynamic_module_v1_type_EventHttpResponseHeadersStatus(result)
}

// __envoy_dynamic_module_v1_event_http_response_body is called when response body data is received.
//
//export __envoy_dynamic_module_v1_event_http_response_body
func __envoy_dynamic_module_v1_event_http_response_body(
	envoyFilterPtr __envoy_dynamic_module_v1_type_EnvoyFilterPtr,
	httpContextPtr __envoy_dynamic_module_v1_type_HttpContextPtr,
	buffer __envoy_dynamic_module_v1_type_HttpResponseBodyBufferPtr,
	endOfStream __envoy_dynamic_module_v1_type_EndOfStream) __envoy_dynamic_module_v1_type_EventHttpResponseBodyStatus {
	envoyPtr := EnvoyFilterPtr{Raw: envoyFilterPtr}
	httpCtx := unwrapRawPinHttpContext(uintptr(httpContextPtr))
	buf := ResponseBodyBufferPtr{Raw: buffer}
	end := endOfStream != 0
	result := httpCtx.EventHttpResponseBody(envoyPtr, buf, end)
	return __envoy_dynamic_module_v1_type_EventHttpResponseBodyStatus(result)
}

// __envoy_dynamic_module_v1_event_http_destroy is called when the stream is destroyed.
//
//export __envoy_dynamic_module_v1_event_http_destroy
func __envoy_dynamic_module_v1_event_http_destroy(
	envoyFilterPtr __envoy_dynamic_module_v1_type_EnvoyFilterPtr,
	httpContextPtr __envoy_dynamic_module_v1_type_HttpContextPtr) {
	envoyPtr := EnvoyFilterPtr{Raw: envoyFilterPtr}
	httpCtx := unwrapRawPinHttpContext(uintptr(httpContextPtr))
	httpCtx.EventHttpDestroy(envoyPtr)
	memManager.removeHttpContext((*pinedHttpContext)(unsafe.Pointer(httpContextPtr)))
}
