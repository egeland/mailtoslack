FROM golang:1-alpine

MAINTAINER Frode Egeland <egeland@gmail.com>

WORKDIR /go/src/app

COPY mailtoslack.go app.go

RUN apk add --virtual buildstuff git --no-cache && \
    go get && \
    go build -i && \
    mv app /usr/local/bin/ && \
    cd / && rm -rf /go && \
    apk del buildstuff

WORKDIR /usr/local/bin

ENV PORT=2525
EXPOSE 2525

CMD ["/usr/local/bin/app"]
