FROM golang:1.21-alpine as dev 

WORKDIR /app

FROM golang:1.21-alpine as build

WORKDIR /app


RUN go build -o gen-blog
