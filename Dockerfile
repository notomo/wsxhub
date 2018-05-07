FROM golang:alpine3.7
RUN apk add --no-cache git=2.15.0-r1 \
    && go get github.com/notomo/wsxhub/cmd/wsxhubd \
    && apk del git
ENTRYPOINT ["wsxhubd"]
