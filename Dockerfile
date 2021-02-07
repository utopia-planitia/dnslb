
FROM golang:1.14.3-alpine3.11 as builder

WORKDIR /go/src/dnslb
COPY go.mod /go/src/dnslb
COPY go.sum /go/src/dnslb

RUN go mod download
RUN go mod verify

COPY cmd /go/src/dnslb/cmd
COPY pkg /go/src/dnslb/pkg
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s -extldflags '-static'" -o /go/bin/dnslb ./cmd/dnslb

FROM gcr.io/distroless/static@sha256:3b3507aa7ce16b085215915d6334f140d8f50ef35bbf3251ee52d3102296f60e

COPY --from=builder /go/bin/dnslb /dnslb
ENTRYPOINT ["/dnslb"]
