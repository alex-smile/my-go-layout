FROM golang:1.21.9 AS builder

COPY ./ /app
WORKDIR /app

ARG BINARY=mygo
ARG VERSION=latest
ARG COMMIT=HEAD
RUN make build -e VERSION=${VERSION} -e COMMIT=${COMMIT} && chmod +x ${BINARY}

# install dlv
RUN go install github.com/go-delve/delve/cmd/dlv@v1.21.0

FROM debian:bullseye-slim

ARG BINARY=mygo

RUN mkdir -p /app/logs
COPY --from=builder /app/${BINARY} /app/${BINARY}
COPY --from=builder /go/bin/dlv /usr/local/bin/dlv

CMD ["/app/mygo", "-c", "/app/config.yaml"]
