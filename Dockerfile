FROM golang:1.17-bullseye

WORKDIR /app

COPY . .

RUN go install
