# syntax=docker/dockerfile:1
#
# Container image for the go-persian-tools demo (the examples/ program).
# Multi-arch via cross-compilation — the build stage always runs on the
# builder's native platform and cross-compiles for the target.

FROM --platform=$BUILDPLATFORM golang:1.23-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG TARGETOS TARGETARCH
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath -ldflags="-s -w" -o /out/persian-tools ./examples

FROM gcr.io/distroless/static-debian12:nonroot
LABEL org.opencontainers.image.source="https://github.com/amiranmanesh/go-persian-tools" \
      org.opencontainers.image.description="Demo runner for the go-persian-tools library" \
      org.opencontainers.image.licenses="MIT"
COPY --from=build /out/persian-tools /persian-tools
ENTRYPOINT ["/persian-tools"]
