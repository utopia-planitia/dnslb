
FROM golang:1.19.3-alpine@sha256:8bd8a4b55b233ea77a81250f83637553ef9e3348c5a0cc3ce440c5d45bf7d51d as builder

WORKDIR /go/src/dnslb
COPY go.mod /go/src/dnslb
COPY go.sum /go/src/dnslb

RUN go mod download
RUN go mod verify

COPY cmd /go/src/dnslb/cmd
COPY pkg /go/src/dnslb/pkg
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s -extldflags '-static'" -o /go/bin/dnslb ./cmd/dnslb

FROM gcr.io/distroless/static@sha256:5759d194607e472ff80fff5833442d3991dd89b219c96552837a2c8f74058617

COPY --from=builder /go/bin/dnslb /dnslb
ENTRYPOINT ["/dnslb"]
