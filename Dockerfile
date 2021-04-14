
FROM golang:1.14.15-alpine3.11@sha256:4f1c80d88c5879067f063770c774a8ffd4de47b684333cdbe9a4ce661931b9b8 as builder

WORKDIR /go/src/dnslb
COPY go.mod /go/src/dnslb
COPY go.sum /go/src/dnslb

RUN go mod download
RUN go mod verify

COPY cmd /go/src/dnslb/cmd
COPY pkg /go/src/dnslb/pkg
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s -extldflags '-static'" -o /go/bin/dnslb ./cmd/dnslb

FROM gcr.io/distroless/static@sha256:60a7d0c45932b6152b2f7ba561db2f91f58ab14aa90b895c58f72062c768fd77

COPY --from=builder /go/bin/dnslb /dnslb
ENTRYPOINT ["/dnslb"]
