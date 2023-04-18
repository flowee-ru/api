FROM golang:1.20

WORKDIR /app

# download modules
COPY go.mod go.sum ./
RUN go mod download

# copy files
COPY . .

# build
RUN CGO_ENABLED=0 GOOS=linux go build -o flowee_api

EXPOSE 8000

# run
CMD ["flowee_api"]