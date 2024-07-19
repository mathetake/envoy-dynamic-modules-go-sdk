package envoy

import "C"

// This file corresponds to the types defined in https://github.com/envoyproxyx/abi/blob/main/abi.h.

// __envoy_dynamic_module_v1_type_ModuleConfigPtr is a pointer to the configuration passed to the
// __envoy_dynamic_module_v1_event_module_init function. Envoy owns the memory of the configuration
// and the module is not supposed to take ownership of it.
type __envoy_dynamic_module_v1_type_ModuleConfigPtr *byte

// __envoy_dynamic_module_v1_type_ModuleConfigSize is the size of the configuration passed to the
// __envoy_dynamic_module_v1_event_module_init function.
type __envoy_dynamic_module_v1_type_ModuleConfigSize int

// __envoy_dynamic_module_v1_type_ModuleContextPtr is a pointer to in-module singleton context
// corresponding to the module. This is passed to __envoy_dynamic_module_v1_event_http_context_init.
type __envoy_dynamic_module_v1_type_ModuleContextPtr uintptr

// __envoy_dynamic_module_v1_type_EnvoyFilterPtr is a pointer to the DynamicModule::HttpFilter
// instance. It is always passed to the module's event hooks. Modules are not supposed to manipulate
// this pointer.
type __envoy_dynamic_module_v1_type_EnvoyFilterPtr uintptr

// __envoy_dynamic_module_v1_type_HttpContextPtr is a pointer to in-module context corresponding
// to a single DynamicModule::HttpFilter instance. It is always passed to the module's event hooks.
type __envoy_dynamic_module_v1_type_HttpContextPtr uintptr

// __envoy_dynamic_module_v1_type_HttpRequestHeadersMapPtr is a pointer to the header map instance.
// This is passed to the __envoy_dynamic_module_v1_event_http_request_headers event hook.
// Modules are not supposed to manipulate this pointer.
type __envoy_dynamic_module_v1_type_HttpRequestHeadersMapPtr uintptr

// __envoy_dynamic_module_v1_type_EventHttpRequestHeadersStatus is the return value of the
// __envoy_dynamic_module_v1_event_http_request_headers event. It should be one of the values
// defined in the FilterHeadersStatus enum.
type __envoy_dynamic_module_v1_type_EventHttpRequestHeadersStatus int

// __envoy_dynamic_module_v1_type_HttpResponseHeaderMapPtr is a pointer to the header map instance.
// This is passed to the __envoy_dynamic_module_v1_event_http_response_headers event hook.
// Modules are not supposed to manipulate this pointer.
type __envoy_dynamic_module_v1_type_HttpResponseHeaderMapPtr uintptr

// __envoy_dynamic_module_v1_type_EventHttpResponseHeadersStatus is the return value of the
// __envoy_dynamic_module_v1_event_http_response_headers event. It should be one of the values
// defined in the FilterHeadersStatus enum.
type __envoy_dynamic_module_v1_type_EventHttpResponseHeadersStatus int

// __envoy_dynamic_module_v1_type_HttpRequestBodyBufferPtr is a pointer to the body buffer instance
// passed via __envoy_dynamic_module_v1_event_http_request_body event hook.
// Modules are not supposed to manipulate this pointer directly.
type __envoy_dynamic_module_v1_type_HttpRequestBodyBufferPtr uintptr

// __envoy_dynamic_module_v1_type_EventHttpRequestBodyStatus is the return value of the
// __envoy_dynamic_module_v1_event_http_request_body event. It should be one of the values defined
// in the FilterDataStatus enum.
type __envoy_dynamic_module_v1_type_EventHttpRequestBodyStatus int

// __envoy_dynamic_module_v1_type_HttpResponseBodyBufferPtr is a pointer to the body buffer instance
// passed via __envoy_dynamic_module_v1_event_http_response_body event hook.
// Modules are not supposed to manipulate this pointer directly.
type __envoy_dynamic_module_v1_type_HttpResponseBodyBufferPtr uintptr

// __envoy_dynamic_module_v1_type_EventHttpResponseBodyStatus is the return value of the
// __envoy_dynamic_module_v1_event_http_response_body event. It should be one of the values defined
// in the FilterDataStatus enum.
type __envoy_dynamic_module_v1_type_EventHttpResponseBodyStatus int

// __envoy_dynamic_module_v1_type_EndOfStream is a boolean value indicating whether the stream has
// reached the end. The value should be 0 if the stream has not reached the end, and 1 if the stream
// has reached the end.
type __envoy_dynamic_module_v1_type_EndOfStream int

// __envoy_dynamic_module_v1_type_InModuleBufferPtr is a pointer to a buffer that is managed by the
// module.
type __envoy_dynamic_module_v1_type_InModuleBufferPtr *byte

// __envoy_dynamic_module_v1_type_InModuleBufferLength is the length of the buffer.
type __envoy_dynamic_module_v1_type_InModuleBufferLength int

// __envoy_dynamic_module_v1_type_DataSlicePtr is a pointer to a buffer that is managed by Envoy.
// This is used to pass buffer slices to the module.
type __envoy_dynamic_module_v1_type_DataSlicePtr *byte

// __envoy_dynamic_module_v1_type_DataSliceLength is the length of the buffer slice.
type __envoy_dynamic_module_v1_type_DataSliceLength int
