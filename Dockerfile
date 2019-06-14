FROM alpine:3.9

RUN apk update && apk add --no-cache ca-certificates

ADD drone-ssh /bin/
ENTRYPOINT ["/bin/drone-ssh"]