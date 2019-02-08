FROM golang:1.11-alpine as build

WORKDIR /app

RUN apk --no-cache add git

RUN go get gopkg.in/urfave/cli.v1

COPY main.go .

RUN go build -o confs.tech.push

FROM alpine:3.9

WORKDIR /app
VOLUME [ "/app/state.json" ]

RUN apk add --no-cache ca-certificates

COPY --from=build /app/confs.tech.push .

ENTRYPOINT [ "./confs.tech.push" ]
