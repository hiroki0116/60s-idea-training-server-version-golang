# Development
FROM golang:alpine as builder

WORKDIR /app
COPY go.mod .
COPY go.sum .

RUN go mod download
RUN go install github.com/githubnemo/CompileDaemon
COPY . .

RUN GOOS=linux GOARCH=amd64 go build -o /main
EXPOSE 8080
ENTRYPOINT CompileDaemon --build="go build main.go" --command=./main

# Production
FROM alpine:3.9

COPY --from=builder /main .

ENV PORT=${PORT}
ENTRYPOINT ["/main"]