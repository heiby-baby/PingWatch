FROM golang:1.19-alpine

WORKDIR /app

COPY . .

RUN go mod init pinger && go mod tidy

RUN go build -o pinger .

CMD ["./pinger"]