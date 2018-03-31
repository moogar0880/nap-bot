# Build Stage
FROM golang:alpine as golang
WORKDIR /go/src/github.com/moogar0880/nap-bot
COPY . /go/src/github.com/moogar0880/nap-bot
RUN apk add --update --no-cache \
        gcc \
        git \
        make \
        musl-dev \
 && make build/alpine

# essentials stage
FROM alpine:latest as alpine
RUN apk --no-cache add tzdata zip ca-certificates
WORKDIR /usr/share/zoneinfo
# -0 means no compression.  Needed because go's tz loader doesn't handle compressed data.
RUN zip -r -0 /zoneinfo.zip .

# production image stage
FROM scratch
COPY --from=golang /go/src/github.com/moogar0880/nap-bot/bin /bin

ENV ZONEINFO /zoneinfo.zip
COPY --from=alpine /zoneinfo.zip /

COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/bin"]
