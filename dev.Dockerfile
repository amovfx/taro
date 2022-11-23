# IMAGE FOR BUILDING
FROM golang:1.19.2 as builder 

WORKDIR /app

COPY . /app

ENV CGO_ENABLED=0

RUN make install

# FINAL IMAGE
FROM alpine as final

COPY --from=builder /go/bin/tarod /bin/
COPY --from=builder /go/bin/tarocli /bin/

COPY ./docker-entrypoint.sh /entrypoint.sh

RUN chmod a+x /entrypoint.sh
VOLUME ["/home/taro/.taro"]

EXPOSE 10029
EXPOSE 8089

ENTRYPOINT ["/entrypoint.sh"]