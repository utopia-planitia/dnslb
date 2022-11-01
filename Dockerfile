
FROM golang:1.19.3-alpine@sha256:8558ae624304387d18694b9ea065cc9813dd4f7f9bd5073edb237541f2d0561b as builder

WORKDIR /go/src/dnslb
COPY go.mod /go/src/dnslb
COPY go.sum /go/src/dnslb

RUN go mod download
RUN go mod verify

COPY cmd /go/src/dnslb/cmd
COPY pkg /go/src/dnslb/pkg
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s -extldflags '-static'" -o /go/bin/dnslb ./cmd/dnslb

FROM gcr.io/distroless/static@sha256:cb0f70353c21d1e472a9ed2055c8a62fd842de52dd2db9dfb31c6e019648bf51

COPY --from=builder /go/bin/dnslb /dnslb
ENTRYPOINT ["/dnslb"]
