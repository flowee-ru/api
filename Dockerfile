FROM golang:1.20

WORKDIR /app

# copy files
COPY . .

# download modules
RUN go mod download

# build
RUN CGO_ENABLED=0 GOOS=linux go build -o flowee_api

# run
CMD ["./flowee_api"]