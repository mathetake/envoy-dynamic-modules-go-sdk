# EnvoyX Go SDK

This is the Go SDK for the EnvoyX modules. The modules are shared libraries that can be loaded by the EnvoyX proxy to extend HTTP filtering capabilities.

The shared library must be compiled with the same environment as EnvoyX, that means the programs must be compiled
on amd64 Linux with the same version of glibc as the EnvoyX proxy.

Since only one Go-based shared library can exist in a process due to [the limitation of the Go runtime](https://github.com/golang/go/issues/65050),
this SDK facilitates the creation of Go-based shared libraries that can be loaded at multiple HTTP filter chain
in Envoy configuration. See the [example](./example) for more details.

## On an amd64 Linux machine

To install the EnvoyX binary locally, the easiest way is to copy the binary from the Docker container:
```bash
docker run --entrypoint=/bin/bash --rm -v $(pwd):/work/envoyx -w /work/envoyx ghcr.io/envoyproxyx/envoy:v1.30-latest-envoyx-main -c "cp /usr/local/bin/envoy /work/envoyx/envoy-bin"
mv envoy-bin /usr/local/bin/envoy
```

where `v1.30` is the Envoy version, and `main` is the [envoyproxyx/envoy](https://github.com/envoyproxyx/envoyx) repository's version (main or tags).
See [github/workflows/commit.yaml](.github/workflows/commit.yaml) for the currently supported versions.

You can build the example and run tests with the following commands:

```bash
make build
make test
```
and if the test passes, you can assume the shared library is compatible with the EnvoyX.

## Others

You can use the Dockerfile provided in this repository to build the shared library.
```
docker build . --platform linux/amd64 -t envoyx-go-sdk:latest
docker run --platform linux/amd64 -v $(pwd):/work/gosdk -w /work/gosdk -it envoyx-go-sdk:latest
```
then follow the instructions in the Linux section.
