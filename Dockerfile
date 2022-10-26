
FROM golang:1.19.2-alpine@sha256:845f16d6c1c1501505a9f35978494bcd77a03f4f0cfeef56e3d8788325bea4a3 as builder

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
