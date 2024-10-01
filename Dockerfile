FROM golang:1.23.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o tiny-bank .

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/tiny-bank .

EXPOSE 8080

CMD [ "./tiny-bank" ]