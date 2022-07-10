FROM golang:1.17-alpine3.15 AS builder

RUN apk add --no-cache git gcc musl-dev make

WORKDIR /go/src/github.com/migratemgr8/mgr8

ENV GO111MODULE=on

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN make build

FROM alpine:3.15

COPY --from=builder /go/src/github.com/migratemgr8/mgr8/bin/mgr8 /usr/local/bin/mgr8
RUN ln -s /usr/local/bin/mgr8 /mgr8

ENTRYPOINT ["mgr8"]
CMD ["--help"]