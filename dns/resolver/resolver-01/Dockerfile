FROM docker.io/library/alpine:latest

RUN apk add --no-cache unbound

EXPOSE 53/udp 53/tcp

CMD ["unbound", "-d"]
