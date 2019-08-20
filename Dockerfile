FROM golang:1.11-alpine as build

RUN apk --no-cache add git

RUN mkdir -p /go/src/github.com/flix-tech/confs.tech.push
WORKDIR /go/src/github.com/flix-tech/confs.tech.push

RUN go get gopkg.in/urfave/cli.v1 github.com/gorilla/feeds

COPY cmd cmd/
COPY confs confs/
COPY main.go .

RUN mkdir /app
RUN go build -o /app/confs.tech.push

FROM alpine:3.9

WORKDIR /app
VOLUME [ "/app/state.json" ]

RUN apk add --no-cache ca-certificates

COPY --from=build /app/confs.tech.push .

ENTRYPOINT [ "./confs.tech.push" ]
