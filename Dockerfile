
FROM golang:1.22.1-alpine@sha256:fc5e5848529786cf1136563452b33d713d5c60b2c787f6b2a077fa6eeefd9114 as builder

WORKDIR /go/src/dnslb
COPY go.mod /go/src/dnslb
COPY go.sum /go/src/dnslb

RUN go mod download
RUN go mod verify

COPY cmd /go/src/dnslb/cmd
COPY pkg /go/src/dnslb/pkg
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s -extldflags '-static'" -o /go/bin/dnslb ./cmd/dnslb

FROM gcr.io/distroless/static@sha256:7e5c6a2a4ae854242874d36171b31d26e0539c98fc6080f942f16b03e82851ab

COPY --from=builder /go/bin/dnslb /dnslb
ENTRYPOINT ["/dnslb"]
