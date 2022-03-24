FROM golang:alpine

RUN apk add build-base
RUN mkdir /app
WORKDIR /app

COPY . /app

RUN go mod tidy
RUN go mod download

RUN go build -o /gobot cmd/main.go

EXPOSE 8080

ENTRYPOINT ["/gobot"]
