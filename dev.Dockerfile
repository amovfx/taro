# IMAGE FOR BUILDING
FROM golang:1.19.3-alpine3.16

WORKDIR /app

COPY . /app

ENV CGO_ENABLED=0

RUN apk --no-cache add \
    bash \
    su-exec \
    ca-certificates \
    make

# FINAL IMAGE

RUN make install
RUN mv /go/bin/tarod /bin/
RUN mv /go/bin/tarocli /bin/


COPY docker-entrypoint.sh /entrypoint.sh


RUN chmod a+x /entrypoint.sh
VOLUME ["/home/taro/.taro"]

EXPOSE 10029
EXPOSE 8089

ENTRYPOINT ["/entrypoint.sh"]

CMD ["tarod"]
