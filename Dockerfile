FROM --platform=$BUILDPLATFORM golang:1.21-alpine AS builder

WORKDIR /go/src/github.com/stackvista/sts-otel-bridge
COPY . .

ARG TARGETOS TARGETARCH
ENV GOOS $TARGETOS
ENV GOARCH $TARGETARCH
RUN go mod tidy
RUN go build -o sts-otel-bridge


FROM gcr.io/distroless/static-debian11
ENTRYPOINT ["/usr/bin/sts-otel-bridge"]
COPY --from=builder /go/src/github.com/stackvista/sts-otel-bridge/sts-otel-bridge /usr/bin/sts-otel-bridge
