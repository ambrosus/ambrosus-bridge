FROM golang:1.18-alpine
RUN apk add build-base

WORKDIR /app
COPY ./go.mod ./go.sum ./
RUN go mod download

COPY . .

RUN go build -o bridge cmd/bridge/main.go
EXPOSE 9090
EXPOSE 8080

CMD ["/app/bridge"]
