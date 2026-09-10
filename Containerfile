FROM docker.io/library/golang:latest


WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o ensor .

CMD ["./ensor"]
