# EnvoyX Go SDK example

This example demonstrates how to use the EnvoyX Go SDK to create HTTP filters for Envoy. This example is supposed to be compiled as a
single shared library but to server multiple HTTP filters. See [envoy.yaml](envoy.yaml) for the configuration.

In main.go, this multiplexes the different HTTP filters based on the `filter_config` parameter given in the Envoy configuration.
Each file named `filter_<name>.go` is a separate HTTP filter implementation which is run on the separater HTTP filter chain.

Please note that many examples are also designed as an E2E test for this SDK itself.

## Build and Run

To build and run the example, run the following, assuming you have `envoy` installed as exaplined in [README.md](../README.md):

```bash

```
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -buildmode=c-shared -o main .
envoy -c envoy.yaml
```