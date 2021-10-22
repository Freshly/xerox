FROM golang:1.17-bullseye
RUN apt-get update \
    && apt-get install -y \
    redis-server
WORKDIR /app
COPY . .
RUN go install
