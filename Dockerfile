
FROM golang:1.20.7-alpine@sha256:b255d9352ff4967f0424731acf5a262c3bb4b5bd712e6324ad9ed533c9d2aee7 as builder

WORKDIR /go/src/dnslb
COPY go.mod /go/src/dnslb
COPY go.sum /go/src/dnslb

RUN go mod download
RUN go mod verify

COPY cmd /go/src/dnslb/cmd
COPY pkg /go/src/dnslb/pkg
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s -extldflags '-static'" -o /go/bin/dnslb ./cmd/dnslb

FROM gcr.io/distroless/static@sha256:7198a357ff3a8ef750b041324873960cf2153c11cc50abb9d8d5f8bb089f6b4e

COPY --from=builder /go/bin/dnslb /dnslb
ENTRYPOINT ["/dnslb"]
