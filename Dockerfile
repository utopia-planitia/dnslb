
FROM golang:1.20.2-alpine@sha256:1db127655b32aa559e32ed3754ed2ea735204d967a433e4b605aed1dd44c5084 as builder

WORKDIR /go/src/dnslb
COPY go.mod /go/src/dnslb
COPY go.sum /go/src/dnslb

RUN go mod download
RUN go mod verify

COPY cmd /go/src/dnslb/cmd
COPY pkg /go/src/dnslb/pkg
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s -extldflags '-static'" -o /go/bin/dnslb ./cmd/dnslb

FROM gcr.io/distroless/static@sha256:97b762efb017cbbabf566046852de8049f84f73e168282d06da316851c7ef263

COPY --from=builder /go/bin/dnslb /dnslb
ENTRYPOINT ["/dnslb"]
