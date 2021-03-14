FROM golang:alpine
RUN apk update && \
    apk upgrade && \
    rm -rf /var/cache/apk/*
WORKDIR /go/src/app
COPY . .

RUN go build -o bin/erised -ldflags "-s -w"

EXPOSE 8080

CMD ["bin/erised"]
