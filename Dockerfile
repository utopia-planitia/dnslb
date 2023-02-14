
FROM golang:1.20.0-alpine@sha256:0d145ecb3cb3772ee54d3a97ae2774aa4f8a179f28f9d4ea67b9cb38b58acebd as builder

WORKDIR /go/src/dnslb
COPY go.mod /go/src/dnslb
COPY go.sum /go/src/dnslb

RUN go mod download
RUN go mod verify

COPY cmd /go/src/dnslb/cmd
COPY pkg /go/src/dnslb/pkg
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s -extldflags '-static'" -o /go/bin/dnslb ./cmd/dnslb

FROM gcr.io/distroless/static@sha256:0f29b9f9818cc0371cae6cef4723f0406184547deeb3aee0db831aeea4527605

COPY --from=builder /go/bin/dnslb /dnslb
ENTRYPOINT ["/dnslb"]
