FROM golang:1.14

WORKDIR /app/golang-rest-api

ADD . .

RUN go get -d -v ./...

RUN go build -o ./build/golang-rest-api .

EXPOSE 8080

CMD ["./build/golang-rest-api"]