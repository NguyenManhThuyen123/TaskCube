FROM golang:1.17.3-alpine3.14

WORKDIR /app

RUN go mod init app

COPY . .
RUN go mod tidy

EXPOSE 8080
CMD ["go", "run", "."]