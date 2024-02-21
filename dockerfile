FROM golang:1.21-alpine as build

WORKDIR /app

COPY app/* .

RUN go build -o gen-blog .

FROM alpine as runtime 

COPY --from=build /app/gen-blog usr/local/bin/gen-blog

COPY run.sh /

ENTRYPOINT [ "./run.sh" ]

