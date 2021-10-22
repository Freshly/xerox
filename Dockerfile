FROM golang:1.17-bullseye
RUN apt-get install redis-tools
WORKDIR /app
COPY . .
RUN go install
