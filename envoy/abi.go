//go:build cgo

package envoy

/*
#include "abi.h"
*/
import "C"
import (
	"runtime"
	"unsafe"
)

//export __envoy_dynamic_module_v1_event_program_init
func __envoy_dynamic_module_v1_event_program_init() C.size_t {
	return 0
}

//export __envoy_dynamic_module_v1_event_http_filter_init
func __envoy_dynamic_module_v1_event_http_filter_init(
	configPtr C.__envoy_dynamic_module_v1_type_HttpFilterConfigPtr,
	configSize C.__envoy_dynamic_module_v1_type_HttpFilterConfigSize) C.__envoy_dynamic_module_v1_type_HttpFilterPtr {
	rawStr := unsafe.String((*byte)(unsafe.Pointer(uintptr(configPtr))), configSize)
	// Copy the config string to Go memory so that the caller can take ownership of the memory.
	var configStrCopy = make([]byte, len(rawStr))
	copy(configStrCopy, rawStr)
	// Call the exported function from the Go module.
	httpFilter := NewHttpFilter(rawStr)
	pined := memManager.pinHttpFilter(httpFilter)
	return C.__envoy_dynamic_module_v1_type_HttpFilterPtr((uintptr)(unsafe.Pointer(pined)))
}

//export __envoy_dynamic_module_v1_event_http_filter_destroy
func __envoy_dynamic_module_v1_event_http_filter_destroy(
	httpFilterPtr C.__envoy_dynamic_module_v1_type_HttpFilterPtr) {
	httpFilter := memManager.unwrapPinnedHttpFilter(uintptr(httpFilterPtr))
	httpFilter.filter.Destroy()
	memManager.unpinHttpFilter(httpFilter)
}

//export __envoy_dynamic_module_v1_event_http_filter_instance_init
func __envoy_dynamic_module_v1_event_http_filter_instance_init(
	envoyFilterPtr C.__envoy_dynamic_module_v1_type_EnvoyFilterInstancePtr,
	moduleCtx C.__envoy_dynamic_module_v1_type_HttpFilterPtr,
) C.__envoy_dynamic_module_v1_type_HttpFilterInstancePtr {
	envoyPtr := &envoyFilterC{raw: envoyFilterPtr}
	m := *(*HttpFilter)(unsafe.Pointer(uintptr(moduleCtx))) //nolint:govet
	httpInstance := m.NewHttpFilterInstance(envoyPtr)
	pined := memManager.pinHttpFilterInstance(httpInstance)
	pined.envoyFilter = envoyPtr
	return C.__envoy_dynamic_module_v1_type_HttpFilterInstancePtr(uintptr((unsafe.Pointer(pined))))
}

//export __envoy_dynamic_module_v1_event_http_filter_instance_request_headers
func __envoy_dynamic_module_v1_event_http_filter_instance_request_headers(
	httpFilterInstancePtr C.__envoy_dynamic_module_v1_type_HttpFilterInstancePtr,
	requestHeadersPtr C.__envoy_dynamic_module_v1_type_HttpRequestHeadersMapPtr,
	endOfStream C.__envoy_dynamic_module_v1_type_EndOfStream,
) C.__envoy_dynamic_module_v1_type_EventHttpRequestHeadersStatus {
	httpInstance := unwrapRawPinHttpFilterInstance(uintptr(httpFilterInstancePtr))
	mapPtr := requestHeadersC{raw: requestHeadersPtr}
	end := endOfStream != 0
	result := httpInstance.filterInstance.EventHttpRequestHeaders(mapPtr, end)
	return C.__envoy_dynamic_module_v1_type_EventHttpRequestHeadersStatus(result)
}

//export __envoy_dynamic_module_v1_event_http_filter_instance_request_body
func __envoy_dynamic_module_v1_event_http_filter_instance_request_body(
	httpFilterInstancePtr C.__envoy_dynamic_module_v1_type_HttpFilterInstancePtr,
	buffer C.__envoy_dynamic_module_v1_type_HttpRequestBodyBufferPtr,
	endOfStream C.__envoy_dynamic_module_v1_type_EndOfStream) C.__envoy_dynamic_module_v1_type_EventHttpRequestBodyStatus {
	httpInstance := unwrapRawPinHttpFilterInstance(uintptr(httpFilterInstancePtr))
	buf := requestBodyBufferC{raw: buffer}
	end := endOfStream != 0
	result := httpInstance.filterInstance.EventHttpRequestBody(buf, end)
	return C.__envoy_dynamic_module_v1_type_EventHttpRequestBodyStatus(result)
}

//export __envoy_dynamic_module_v1_event_http_filter_instance_response_headers
func __envoy_dynamic_module_v1_event_http_filter_instance_response_headers(
	httpFilterInstancePtr C.__envoy_dynamic_module_v1_type_HttpFilterInstancePtr,
	responseHeadersMapPtr C.__envoy_dynamic_module_v1_type_HttpResponseHeaderMapPtr,
	endOfStream C.__envoy_dynamic_module_v1_type_EndOfStream) C.__envoy_dynamic_module_v1_type_EventHttpResponseHeadersStatus {
	httpInstance := unwrapRawPinHttpFilterInstance(uintptr(httpFilterInstancePtr))
	mapPtr := responseHeadersC{raw: responseHeadersMapPtr}
	end := endOfStream != 0
	result := httpInstance.filterInstance.EventHttpResponseHeaders(mapPtr, end)
	return C.__envoy_dynamic_module_v1_type_EventHttpResponseHeadersStatus(result)
}

//export __envoy_dynamic_module_v1_event_http_filter_instance_response_body
func __envoy_dynamic_module_v1_event_http_filter_instance_response_body(
	httpFilterInstancePtr C.__envoy_dynamic_module_v1_type_HttpFilterInstancePtr,
	buffer C.__envoy_dynamic_module_v1_type_HttpResponseBodyBufferPtr,
	endOfStream C.__envoy_dynamic_module_v1_type_EndOfStream) C.__envoy_dynamic_module_v1_type_EventHttpResponseBodyStatus {
	httpInstance := unwrapRawPinHttpFilterInstance(uintptr(httpFilterInstancePtr))
	buf := responseBodyBufferC{raw: buffer}
	end := endOfStream != 0
	result := httpInstance.filterInstance.EventHttpResponseBody(buf, end)
	return C.__envoy_dynamic_module_v1_type_EventHttpResponseBodyStatus(result)
}

//export __envoy_dynamic_module_v1_event_http_filter_instance_destroy
func __envoy_dynamic_module_v1_event_http_filter_instance_destroy(
	httpFilterInstancePtr C.__envoy_dynamic_module_v1_type_HttpFilterInstancePtr) {
	httpInstance := unwrapRawPinHttpFilterInstance(uintptr(httpFilterInstancePtr))
	httpInstance.filterInstance.EventHttpDestroy(httpInstance.envoyFilter)
	httpInstance.envoyFilter.(*envoyFilterC).destroyed = true
	memManager.unpinHttpFilterInstance((*pinedHttpFilterInstance)(unsafe.Pointer(uintptr(httpFilterInstancePtr))))
}

// envoyFilterC implements the EnvoyFilterInstance.
type envoyFilterC struct {
	raw       C.__envoy_dynamic_module_v1_type_EnvoyFilterInstancePtr
	destroyed bool
}

// ContinueRequest implements EnvoyFilterInstance.
func (c *envoyFilterC) ContinueRequest() {
	if c.destroyed {
		return
	}
	C.__envoy_dynamic_module_v1_http_continue_request(c.raw)
}

// ContinueResponse implements EnvoyFilterInstance.
func (c *envoyFilterC) ContinueResponse() {
	if c.destroyed {
		return
	}
	C.__envoy_dynamic_module_v1_http_continue_response(c.raw)
}

// Destroyed implements EnvoyFilterInstance.
func (c *envoyFilterC) Destroyed() bool {
	return c.destroyed
}

// requestHeadersC implements RequestHeaders.
type requestHeadersC struct {
	raw C.__envoy_dynamic_module_v1_type_HttpRequestHeadersMapPtr
}

// responseHeadersC implements ResponseHeaders.
type responseHeadersC struct {
	raw C.__envoy_dynamic_module_v1_type_HttpResponseHeaderMapPtr
}

// requestBodyBufferC implements RequestBodyBuffer.
type requestBodyBufferC struct {
	raw C.__envoy_dynamic_module_v1_type_HttpRequestBodyBufferPtr
}

// responseBodyBufferC implements ResponseBodyBuffer.
type responseBodyBufferC struct {
	raw C.__envoy_dynamic_module_v1_type_HttpResponseBodyBufferPtr
}

var (
	_ RequestHeaders     = (*requestHeadersC)(nil)
	_ ResponseHeaders    = (*responseHeadersC)(nil)
	_ RequestBodyBuffer  = (*requestBodyBufferC)(nil)
	_ ResponseBodyBuffer = (*responseBodyBufferC)(nil)
)

// Get implements RequestHeaders.
func (r requestHeadersC) Get(key string, iter func(value HeaderValue)) {
	// Take the raw pointer to the key by using unsafe.
	keyPtr := uintptr(unsafe.Pointer(unsafe.StringData(key)))
	keySize := len(key)

	var resultPtr *byte
	var resultSize int
	total := C.__envoy_dynamic_module_v1_http_get_request_header_value(r.raw,
		C.__envoy_dynamic_module_v1_type_InModuleBufferPtr(keyPtr),
		C.__envoy_dynamic_module_v1_type_InModuleBufferLength(keySize),
		C.__envoy_dynamic_module_v1_type_DataSliceLengthResult(uintptr(unsafe.Pointer(&resultPtr))),
		C.__envoy_dynamic_module_v1_type_DataSliceLengthResult(uintptr(unsafe.Pointer(&resultSize))),
	)
	if total == 0 {
		return
	}

	iter(HeaderValue{data: resultPtr, size: int(resultSize)})

	for i := 1; i < int(total); i++ {
		C.__envoy_dynamic_module_v1_http_get_request_header_value_nth(r.raw,
			C.__envoy_dynamic_module_v1_type_InModuleBufferPtr(keyPtr),
			C.__envoy_dynamic_module_v1_type_InModuleBufferLength(keySize),
			C.__envoy_dynamic_module_v1_type_DataSliceLengthResult(uintptr(unsafe.Pointer(&resultPtr))),
			C.__envoy_dynamic_module_v1_type_DataSliceLengthResult(uintptr(unsafe.Pointer(&resultSize))),
			C.size_t(i),
		)
		iter(HeaderValue{data: resultPtr, size: int(resultSize)})
	}
	runtime.KeepAlive(key)
}

// Set implements RequestHeaders.
func (r requestHeadersC) Set(key, value string) {
	r.set(
		uintptr(unsafe.Pointer(unsafe.StringData(key))), len(key),
		uintptr(unsafe.Pointer(unsafe.StringData(value))), len(value),
	)
	runtime.KeepAlive(key)
	runtime.KeepAlive(value)
}

// Remove implements RequestHeaders.
func (r requestHeadersC) Remove(key string) {
	r.set(uintptr(unsafe.Pointer(unsafe.StringData(key))), len(key), 0, 0)
	runtime.KeepAlive(key)
}

func (r requestHeadersC) set(keyPtr uintptr, keySize int, valuePtr uintptr, valueSize int) {
	C.__envoy_dynamic_module_v1_http_set_request_header(r.raw,
		C.__envoy_dynamic_module_v1_type_InModuleBufferPtr(keyPtr),
		C.__envoy_dynamic_module_v1_type_InModuleBufferLength(keySize),
		C.__envoy_dynamic_module_v1_type_InModuleBufferPtr(valuePtr),
		C.__envoy_dynamic_module_v1_type_InModuleBufferLength(valueSize),
	)
}

// Get implements ResponseHeaders.
func (r responseHeadersC) Get(key string, iter func(value HeaderValue)) {
	// Take the raw pointer to the key by using unsafe.
	keyPtr := uintptr(unsafe.Pointer(unsafe.StringData(key)))
	keySize := len(key)

	var resultPtr *byte
	var resultSize int
	total := C.__envoy_dynamic_module_v1_http_get_response_header_value(r.raw,
		C.__envoy_dynamic_module_v1_type_InModuleBufferPtr(keyPtr),
		C.__envoy_dynamic_module_v1_type_InModuleBufferLength(keySize),
		C.__envoy_dynamic_module_v1_type_DataSliceLengthResult(uintptr(unsafe.Pointer(&resultPtr))),
		C.__envoy_dynamic_module_v1_type_DataSliceLengthResult(uintptr(unsafe.Pointer(&resultSize))),
	)
	if total == 0 {
		return
	}

	iter(HeaderValue{data: resultPtr, size: resultSize})

	for i := 1; i < int(total); i++ {
		C.__envoy_dynamic_module_v1_http_get_response_header_value_nth(r.raw,
			C.__envoy_dynamic_module_v1_type_InModuleBufferPtr(keyPtr),
			C.__envoy_dynamic_module_v1_type_InModuleBufferLength(keySize),
			C.__envoy_dynamic_module_v1_type_DataSliceLengthResult(uintptr(unsafe.Pointer(&resultPtr))),
			C.__envoy_dynamic_module_v1_type_DataSliceLengthResult(uintptr(unsafe.Pointer(&resultSize))),
			C.size_t(i),
		)
		iter(HeaderValue{data: resultPtr, size: resultSize})
	}

	runtime.KeepAlive(key)
}

// Set implements ResponseHeaders.
func (r responseHeadersC) Set(key, value string) {
	r.set(
		uintptr(unsafe.Pointer(unsafe.StringData(key))), len(key),
		uintptr(unsafe.Pointer(unsafe.StringData(value))), len(value),
	)
	runtime.KeepAlive(key)
	runtime.KeepAlive(value)
}

// Remove implements ResponseHeaders.
func (r responseHeadersC) Remove(key string) {
	r.set(uintptr(unsafe.Pointer(unsafe.StringData(key))), len(key), 0, 0)
	runtime.KeepAlive(key)
}

func (r responseHeadersC) set(keyPtr uintptr, keySize int, valuePtr uintptr, valueSize int) {
	C.__envoy_dynamic_module_v1_http_set_response_header(r.raw,
		C.__envoy_dynamic_module_v1_type_InModuleBufferPtr(keyPtr),
		C.__envoy_dynamic_module_v1_type_InModuleBufferLength(keySize),
		C.__envoy_dynamic_module_v1_type_InModuleBufferPtr(valuePtr),
		C.__envoy_dynamic_module_v1_type_InModuleBufferLength(valueSize),
	)
}
