FROM golang:1.19-alpine

WORKDIR /app

COPY . .

RUN go mod init backend && go mod tidy

RUN go build -o backend .


CMD ["./backend"]