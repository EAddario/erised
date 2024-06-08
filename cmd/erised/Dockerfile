FROM golang:buster AS build
WORKDIR /go/src/app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o /go/bin/erised -ldflags "-s -w"

FROM gcr.io/distroless/base-debian12
COPY --from=build /go/bin/erised /
EXPOSE 8080

ENTRYPOINT ["/erised"]
CMD [""]
