
FROM golang:1.17.6-alpine@sha256:519c827ec22e5cf7417c9ff063ec840a446cdd30681700a16cf42eb43823e27c as builder

WORKDIR /go/src/dnslb
COPY go.mod /go/src/dnslb
COPY go.sum /go/src/dnslb

RUN go mod download
RUN go mod verify

COPY cmd /go/src/dnslb/cmd
COPY pkg /go/src/dnslb/pkg
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s -extldflags '-static'" -o /go/bin/dnslb ./cmd/dnslb

FROM gcr.io/distroless/static@sha256:8ad6f3ec70dad966479b9fb48da991138c72ba969859098ec689d1450c2e6c97

COPY --from=builder /go/bin/dnslb /dnslb
ENTRYPOINT ["/dnslb"]
