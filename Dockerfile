
FROM golang:1.17.1-alpine@sha256:13919fb9091f6667cb375d5fdf016ecd6d3a5d5995603000d422b04583de4ef9 as builder

WORKDIR /go/src/dnslb
COPY go.mod /go/src/dnslb
COPY go.sum /go/src/dnslb

RUN go mod download
RUN go mod verify

COPY cmd /go/src/dnslb/cmd
COPY pkg /go/src/dnslb/pkg
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s -extldflags '-static'" -o /go/bin/dnslb ./cmd/dnslb

FROM gcr.io/distroless/static@sha256:9fc06167562a1ab9af57bd403825899d3b34dc658a2035e3a1ddc67d80f434ef

COPY --from=builder /go/bin/dnslb /dnslb
ENTRYPOINT ["/dnslb"]
