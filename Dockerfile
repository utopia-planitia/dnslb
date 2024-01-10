
FROM golang:1.21.6-alpine@sha256:fd78f2fb1e49bcf343079bbbb851c936a18fc694df993cbddaa24ace0cc724c5 as builder

WORKDIR /go/src/dnslb
COPY go.mod /go/src/dnslb
COPY go.sum /go/src/dnslb

RUN go mod download
RUN go mod verify

COPY cmd /go/src/dnslb/cmd
COPY pkg /go/src/dnslb/pkg
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s -extldflags '-static'" -o /go/bin/dnslb ./cmd/dnslb

FROM gcr.io/distroless/static@sha256:9be3fcc6abeaf985b5ecce59451acbcbb15e7be39472320c538d0d55a0834edc

COPY --from=builder /go/bin/dnslb /dnslb
ENTRYPOINT ["/dnslb"]
