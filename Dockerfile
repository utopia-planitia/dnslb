
FROM golang:1.19.1-alpine@sha256:d475cef843a02575ebdcb1416d98cd76bab90a5ae8bc2cd15f357fc08b6a329f as builder

WORKDIR /go/src/dnslb
COPY go.mod /go/src/dnslb
COPY go.sum /go/src/dnslb

RUN go mod download
RUN go mod verify

COPY cmd /go/src/dnslb/cmd
COPY pkg /go/src/dnslb/pkg
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s -extldflags '-static'" -o /go/bin/dnslb ./cmd/dnslb

FROM gcr.io/distroless/static@sha256:35817ff4e3441ebfa430cfca98121e43cbf4a134412285c2460a9464c97ac2ca

COPY --from=builder /go/bin/dnslb /dnslb
ENTRYPOINT ["/dnslb"]
