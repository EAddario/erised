FROM golang:buster AS build
WORKDIR /go/src/app
COPY . .
RUN go build -o /go/bin/erised -ldflags "-s -w"

FROM gcr.io/distroless/base-debian10
COPY --from=build /go/bin/erised /
EXPOSE 8080

CMD ["/erised"]
