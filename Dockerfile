
FROM golang:1.21.4-alpine@sha256:110b07af87238fbdc5f1df52b00927cf58ce3de358eeeb1854f10a8b5e5e1411 as builder

WORKDIR /go/src/dnslb
COPY go.mod /go/src/dnslb
COPY go.sum /go/src/dnslb

RUN go mod download
RUN go mod verify

COPY cmd /go/src/dnslb/cmd
COPY pkg /go/src/dnslb/pkg
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s -extldflags '-static'" -o /go/bin/dnslb ./cmd/dnslb

FROM gcr.io/distroless/static@sha256:6706c73aae2afaa8201d63cc3dda48753c09bcd6c300762251065c0f7e602b25

COPY --from=builder /go/bin/dnslb /dnslb
ENTRYPOINT ["/dnslb"]
