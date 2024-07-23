package main

import (
	"github.com/envoyproxyx/go-sdk/envoy"
)

func main() {} // main function must be present but empty.

// Set the envoy.NewHttpFilter function to create a new http filter.
func init() { envoy.NewHttpFilter = newHttpFilter }

// newHttpFilter creates a new http filter based on the config.
//
// `config` is the configuration string that is specified in the Envoy configuration.
func newHttpFilter(config string) envoy.HttpFilter {
	switch config {
	case "helloworld":
		return newHelloWorldHttpFilter(config)
	case "delay":
		return newDelayHttpFilter(config)
	case "headers":
		return newHeadersHttpFilter(config)
	default:
		panic("unknown filter: " + config)
	}
}
