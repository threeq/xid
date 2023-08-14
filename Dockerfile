FROM golang:1.21 AS builder
RUN mkdir /app
WORKDIR /app
COPY . .
RUN make build-linux

FROM alpine:3.10
RUN mkdir /app
WORKDIR /app
COPY --from=builder /app/bin/xid_linux ./xid

EXPOSE 8888
EXPOSE 8999

ENTRYPOINT ["/app/xid"]