package envoy

// This file corresponds to the enums defined in abi.h.

// EventHttpRequestHeadersStatus is the return value of the HttpFilterInstance.EventHttpRequestHeaders event.
type EventHttpRequestHeadersStatus int

const (
	// EventHttpRequestHeadersStatusContinue is returned when the operation should continue.
	EventHttpRequestHeadersStatusContinue EventHttpRequestHeadersStatus = 0
	// EventHttpRequestHeadersStatusStopIteration indicates that Envoy shouldn't continue
	// from processing the headers and should stop filter iteration. In other words, HttpFilterInstance.EventHttpRequestBody
	// will be called while not sending headers to the upstream. The header
	// processing can be resumed by either calling EnvoyFilterInstance.ContinueRequest, or returns
	// continue status from the HttpFilterInstance.EventHttpRequestBody.
	EventHttpRequestHeadersStatusStopIteration EventHttpRequestHeadersStatus = 1
	// EventHttpRequestHeadersStatusStopAllIterationAndBuffer indicates
	// that Envoy should stop all iteration and continue to buffer the request body
	// until the limit is reached. When the limit is reached, Envoy will stop buffering and returns 500
	// to the client. This means that HttpFilterInstance.EventHttpRequestBody will not be called.
	//
	// The header processing can be resumed by either calling EnvoyFilterInstance.ContinueRequest, or
	// returns continue status from the HttpFilterInstance.EventHttpRequestBody.
	EventHttpRequestHeadersStatusStopAllIterationAndBuffer EventHttpRequestHeadersStatus = 3
)

// EventHttpRequestBodyStatus is the return value of the HttpFilterInstance.EventHttpRequestBody event.
type EventHttpRequestBodyStatus int

const (
	// EventHttpRequestBodyStatusContinue is returned when the operation should continue.
	EventHttpRequestBodyStatusContinue EventHttpRequestBodyStatus = 0
	// EventHttpRequestBodyStatusStopIterationAndBuffer indicates that Envoy shouldn't continue
	// from processing the body frame and should stop iteration, but continue buffering the body until
	// the limit is reached. When the limit is reached, Envoy will stop buffering and returns 500 to the
	// client.
	//
	// This stops sending body data to the upstream, so if the module wants to continue sending body
	// data, it should call EnvoyFilterInstance.ContinueRequest or return continue status in the
	// subsequent HttpFilterInstance.EventHttpRequestBody calls.
	EventHttpRequestBodyStatusStopIterationAndBuffer EventHttpRequestBodyStatus = 1
)

// EventHttpResponseHeadersStatus is the return value of the HttpFilterInstance.EventHttpResponseHeaders event.
type EventHttpResponseHeadersStatus int

const (
	// EventHttpResponseHeadersStatusContinue is returned when the operation should continue.
	EventHttpResponseHeadersStatusContinue EventHttpResponseHeadersStatus = 0
	// EventHttpResponseHeadersStatusStopIteration indicates that Envoy shouldn't continue
	// from processing the headers and should stop filter iteration. In other words, EventHttpResponseBody
	// will be called while not sending headers to the upstream. The header
	// processing can be resumed by either calling EnvoyFilterInstance.ContinueResponse, or returns
	// continue status from the EventHttpResponseBody.
	EventHttpResponseHeadersStatusStopIteration EventHttpResponseHeadersStatus = 1

	// EventHttpResponseHeadersStatusStopAllIterationAndBuffer indicates
	// that Envoy should stop all iteration and continue to buffer the response body
	// until the limit is reached. When the limit is reached, Envoy will stop buffering and returns 500
	// to the client. This means that HttpFilterInstance.EventHttpResponseBody will not be called.
	//
	// The header processing can be resumed by either calling EnvoyFilterInstance.ContinueResponse, or
	// returns continue status from the HttpFilterInstance.EventHttpResponseBody.
	EventHttpResponseHeadersStatusStopAllIterationAndBuffer EventHttpResponseHeadersStatus = 3
)

// EventHttpResponseBodyStatus is the return value of the HttpFilterInstance.EventHttpResponseBody event.
type EventHttpResponseBodyStatus int

const (
	// EventHttpResponseBodyStatusContinue is returned when the operation should continue.
	EventHttpResponseBodyStatusContinue EventHttpResponseBodyStatus = 0
	// EventHttpResponseBodyStatusStopIterationAndBuffer indicates that Envoy shouldn't continue
	// from processing the body frame and should stop iteration, but continue buffering the body until
	// the limit is reached. When the limit is reached, Envoy will stop buffering and returns 500 to the
	// client.
	//
	// This stops sending body data to the upstream, so if the module wants to continue sending body
	// data, it should call EnvoyFilterInstance.ContinueResponse or return continue status in the
	// subsequent HttpFilterInstance.EventHttpResponseBody calls.
	EventHttpResponseBodyStatusStopIterationAndBuffer EventHttpResponseBodyStatus = 1
)
