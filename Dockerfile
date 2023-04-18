FROM golang:1.20

WORKDIR /app

# download modules
COPY go.mod go.sum ./
RUN go mod download

# copy files
COPY ./**/*.go .

# build
RUN go build -o /flowee_api

EXPOSE 8000

# run
CMD ["/flowee_api"]