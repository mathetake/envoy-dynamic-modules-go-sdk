# EnvoyX Go SDK

This is the Go SDK for the EnvoyX modules. The modules are shared libraries that can be loaded by the EnvoyX proxy to extend HTTP filtering capabilities.

The shared library must be compiled with the same environment as EnvoyX, that means the programs must be compiled
on amd64 Linux with the same version of glibc as the EnvoyX proxy.

## On an amd64 Linux machine

You can build the examples and run tests, which do not require EnvoyX installed, with the following commands:

```bash
make build
make test
```

If you have the EnvoyX binary is installed on the PATH,
you can run the example e2e test with the following command:

```bash
make e2e
```
and if the test passes, you can assume the shared library is compatible with the EnvoyX.

To install the EnvoyX locally, the easiest way is to copy the binary from the Docker container:
```bash
docker run --entrypoint=/bin/bash --rm -v $(pwd):/work/envoyx -w /work/envoyx ghcr.io/envoyproxyx/envoy:v1.30-latest-envoyx-v0.1.0 -c "cp /usr/local/bin/envoy /work/envoyx/envoy-bin"
mv envoy-bin /usr/local/bin/envoy
```

where `v1.30` and `v0.1.0` are the versions of Envoy and EnvoyX respectively.
See [github/workflows/commit.yaml](.github/workflows/commit.yaml) for the currently supported versions.

## Others

You can use the Dockerfile provided in this repository to build the shared library.
```
docker build . --platform linux/amd64 -t envoyx-go-sdk:latest
docker run --platform linux/amd64 -v $(pwd):/work/gosdk -w /work/gosdk -it envoyx-go-sdk:latest
```
then follow the instructions in the Linux section.
