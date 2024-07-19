# The container image is built in the envoyproxyx/envoyx repository.
ARG ENVOY_IMAGE=ghcr.io/envoyproxyx/envoy:v1.30-latest-envoyx-v0.1.0
FROM ${ENVOY_IMAGE} AS ENVOY

FROM golang:1.22 AS BUILDER
COPY --from=ENVOY /usr/local/bin/envoy /usr/local/bin/envoy
