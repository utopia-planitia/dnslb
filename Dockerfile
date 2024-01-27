
FROM golang:1.21.6-alpine@sha256:a6a7f1fcf12f5efa9e04b1e75020931a616cd707f14f62ab5262bfbe109aa84a as builder

WORKDIR /go/src/dnslb
COPY go.mod /go/src/dnslb
COPY go.sum /go/src/dnslb

RUN go mod download
RUN go mod verify

COPY cmd /go/src/dnslb/cmd
COPY pkg /go/src/dnslb/pkg
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s -extldflags '-static'" -o /go/bin/dnslb ./cmd/dnslb

FROM gcr.io/distroless/static@sha256:9be3fcc6abeaf985b5ecce59451acbcbb15e7be39472320c538d0d55a0834edc

COPY --from=builder /go/bin/dnslb /dnslb
ENTRYPOINT ["/dnslb"]
