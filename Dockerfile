FROM golang:1.14

WORKDIR /app/golang-rest-api

#RUN apk add --no-cache git

ADD . .

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o ./build/golang-rest-api .

EXPOSE 8080

CMD ["./build/golang-rest-api"]