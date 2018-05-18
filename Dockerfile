FROM golang:alpine3.7

RUN apk add --no-cache git=2.15.0-r1 \
    && go get -u github.com/golang/dep/cmd/dep

WORKDIR /go/src/github.com/notomo/wsxhub

ENTRYPOINT ["sh", "./run.sh"]
