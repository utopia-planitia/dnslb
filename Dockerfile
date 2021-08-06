
FROM golang:1.16.7-alpine@sha256:3e6a2def9a57f23344a75bd71d9cd79726f0fbaf4b75330be5669773df0e9d4c as builder

WORKDIR /go/src/dnslb
COPY go.mod /go/src/dnslb
COPY go.sum /go/src/dnslb

RUN go mod download
RUN go mod verify

COPY cmd /go/src/dnslb/cmd
COPY pkg /go/src/dnslb/pkg
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s -extldflags '-static'" -o /go/bin/dnslb ./cmd/dnslb

FROM gcr.io/distroless/static@sha256:c9320b754c2fa2cd2dea50993195f104a24f4c7ebe6e0297c6ddb40ce3679e7d

COPY --from=builder /go/bin/dnslb /dnslb
ENTRYPOINT ["/dnslb"]
