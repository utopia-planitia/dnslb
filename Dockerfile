
FROM golang:1.20.2-alpine@sha256:4e6bc0eafc261b6c8ba9bd9999b6698e8cefbe21e6d90fbc10c34599d75dc608 as builder

WORKDIR /go/src/dnslb
COPY go.mod /go/src/dnslb
COPY go.sum /go/src/dnslb

RUN go mod download
RUN go mod verify

COPY cmd /go/src/dnslb/cmd
COPY pkg /go/src/dnslb/pkg
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s -extldflags '-static'" -o /go/bin/dnslb ./cmd/dnslb

FROM gcr.io/distroless/static@sha256:3c5767883bc3ad6c4ad7caf97f313e482f500f2c214f78e452ac1ca932e1bf7f

COPY --from=builder /go/bin/dnslb /dnslb
ENTRYPOINT ["/dnslb"]
