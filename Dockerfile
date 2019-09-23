FROM alpine:3.10

RUN mkdir /app
WORKDIR /app

COPY bin/xid_linux /app/xid

EXPOSE 8080

ENTRYPOINT ["/app/xid"]