FROM golang:1.22-alpine as dev

RUN apk update && apk add git 

workdir /work