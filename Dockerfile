# syntax=docker/dockerfile:1
#
# Container image for the persian-tools CLI.
# Multi-arch via cross-compilation — the build stage always runs on the
# builder's native platform and cross-compiles for the target.

FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS build
WORKDIR /src
# The module has no dependencies, so there is nothing to download: copying the
# sources straight in is both correct and the fastest path.
COPY . .
ARG TARGETOS TARGETARCH VERSION=dev
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath -ldflags="-s -w -X main.version=${VERSION}" \
    -o /out/persian-tools ./cmd/persian-tools

FROM gcr.io/distroless/static-debian12:nonroot
LABEL org.opencontainers.image.source="https://github.com/amiranmanesh/go-persian-tools" \
      org.opencontainers.image.description="Tools for Persian (Iranian) data: text, digits, money and validators" \
      org.opencontainers.image.licenses="MIT"
COPY --from=build /out/persian-tools /persian-tools
ENTRYPOINT ["/persian-tools"]
