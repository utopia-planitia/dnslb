
FROM golang:1.16.6-alpine@sha256:7700a5e57723ea9e9d0bfbfe66e69ae365d7eec53063ba4a1dd93dc71e37034c as builder

WORKDIR /go/src/dnslb
COPY go.mod /go/src/dnslb
COPY go.sum /go/src/dnslb

RUN go mod download
RUN go mod verify

COPY cmd /go/src/dnslb/cmd
COPY pkg /go/src/dnslb/pkg
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s -extldflags '-static'" -o /go/bin/dnslb ./cmd/dnslb

FROM gcr.io/distroless/static@sha256:912bd2c2b9704ead25ba91b631e3849d940f9d533f0c15cf4fc625099ad145b1

COPY --from=builder /go/bin/dnslb /dnslb
ENTRYPOINT ["/dnslb"]
