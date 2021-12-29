FROM golang:1.17-alpine AS build
WORKDIR /go/src/flusso
ENV GO111MODULE on

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -mod=readonly -o /go/bin/flusso ./cmd/flusso

RUN GRPC_HEALTH_PROBE_VERSION=v0.4.6 && \
  wget -qO/go/bin/grpc_health_probe \
  https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
  chmod +x /go/bin/grpc_health_probe

FROM alpine:latest
COPY --from=build /go/bin/flusso /bin/flusso
COPY --from=build /go/bin/grpc_health_probe /bin/grpc_health_probe
ENTRYPOINT ["/bin/flusso"]
