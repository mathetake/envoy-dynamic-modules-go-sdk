# The container image is built in the envoyproxyx/envoyx repository.
ARG ENVOY_IMAGE=ghcr.io/envoyproxyx/envoy:v1.31-latest-envoyx-main
FROM ${ENVOY_IMAGE} AS ENVOY

FROM golang:1.22 AS BUILDER
COPY --from=ENVOY /usr/local/bin/envoy /usr/local/bin/envoy
