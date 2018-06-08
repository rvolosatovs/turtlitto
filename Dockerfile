FROM alpine:3.7
RUN apk --update --no-cache add ca-certificates
ADD ./release/srrs-linux-amd64 /usr/local/bin/srrs
ADD ./release/front /usr/local/srr/front
RUN chmod 755 /usr/local/bin/srrs
ENTRYPOINT ["/usr/local/bin/srrs", "-static", "/usr/local/srr/front", "-unixSocket", "/trc/trc.sock"]
