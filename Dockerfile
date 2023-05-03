FROM golang:1.20-alpine

WORKDIR /usr/src/app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o flowee_api

EXPOSE 8081

CMD ["./flowee_api"]