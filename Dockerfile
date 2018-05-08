FROM alpine:3.7
RUN apk --update --no-cache add ca-certificates
ADD ./release/soccer-robot-remote-linux-amd64 /usr/local/bin/soccer-robot-remote
ADD ./release/front /usr/local/srr/front
RUN chmod 755 /usr/local/bin/soccer-robot-remote
ENTRYPOINT ["/usr/local/bin/soccer-robot-remote", "-static", "/usr/local/srr/front", "-socket", "/trc/trc.sock"]
