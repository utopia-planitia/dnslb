
FROM golang:1.17.0-alpine@sha256:8c7994cfdf4d488799d40d85d83bd41c7fd290e8eed1affc2abd386150750d2d as builder

WORKDIR /go/src/dnslb
COPY go.mod /go/src/dnslb
COPY go.sum /go/src/dnslb

RUN go mod download
RUN go mod verify

COPY cmd /go/src/dnslb/cmd
COPY pkg /go/src/dnslb/pkg
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s -extldflags '-static'" -o /go/bin/dnslb ./cmd/dnslb

FROM gcr.io/distroless/static@sha256:c9320b754c2fa2cd2dea50993195f104a24f4c7ebe6e0297c6ddb40ce3679e7d

COPY --from=builder /go/bin/dnslb /dnslb
ENTRYPOINT ["/dnslb"]
