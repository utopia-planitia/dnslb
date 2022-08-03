
FROM golang:1.19.0-alpine@sha256:f8e128fa8aa891fe29e22e6401686dffef9bd4c3f5b552b09a7c29f7379979c1 as builder

WORKDIR /go/src/dnslb
COPY go.mod /go/src/dnslb
COPY go.sum /go/src/dnslb

RUN go mod download
RUN go mod verify

COPY cmd /go/src/dnslb/cmd
COPY pkg /go/src/dnslb/pkg
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s -extldflags '-static'" -o /go/bin/dnslb ./cmd/dnslb

FROM gcr.io/distroless/static@sha256:21d3f84a4f37c36199fd07ad5544dcafecc17776e3f3628baf9a57c8c0181b3f

COPY --from=builder /go/bin/dnslb /dnslb
ENTRYPOINT ["/dnslb"]
