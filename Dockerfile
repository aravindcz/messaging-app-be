FROM golang:1.21-alpine
WORKDIR /app
COPY . .
RUN go mod init api && go mod tidy
RUN go build -o api
CMD ["./api"]