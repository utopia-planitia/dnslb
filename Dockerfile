
FROM golang:1.17.2-alpine@sha256:5519c8752f6b53fc8818dc46e9fda628c99c4e8fd2d2f1df71e1f184e71f47dc as builder

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
