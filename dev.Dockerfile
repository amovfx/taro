# IMAGE FOR BUILDING
FROM golang:1.18-alpine as builder 

ENV GODEBUG netdns=cgo
RUN apk add --no-cache --update alpine-sdk \
    bash \
    su-exec \
    ca-certificates

WORKDIR /app

COPY . /app

#ENV CGO_ENABLED=0

RUN make install
RUN go install github.com/go-delve/delve/cmd/dlv@latest

# FINAL IMAGE
# FROM alpine as final

# RUN apk add --no-cache --update alpine-sdk \
#     bash \
#     su-exec \
#     ca-certificates


ENV PATH=$PATH:/bin:/go/bin/

COPY docker-entrypoint.sh /entrypoint.sh
VOLUME ["/home/taro/.taro"]
VOLUME ["/home/lnd/.lnd"]

EXPOSE 10029
EXPOSE 8089

RUN chmod +x /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]

CMD ["tarod"]