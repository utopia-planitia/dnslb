
FROM golang:1.17.3-alpine@sha256:7bd9bf011a76c21b70a1e507de7f8d67656e1f293fd72af80aa46f7b054db515 as builder

WORKDIR /go/src/dnslb
COPY go.mod /go/src/dnslb
COPY go.sum /go/src/dnslb

RUN go mod download
RUN go mod verify

COPY cmd /go/src/dnslb/cmd
COPY pkg /go/src/dnslb/pkg
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s -extldflags '-static'" -o /go/bin/dnslb ./cmd/dnslb

FROM gcr.io/distroless/static@sha256:1cc74da80bbf80d89c94e0c7fe22830aa617f47643f2db73f66c8bd5bf510b25

COPY --from=builder /go/bin/dnslb /dnslb
ENTRYPOINT ["/dnslb"]
