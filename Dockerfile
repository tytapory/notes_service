FROM golang:1.23-alpine
WORKDIR /app
COPY src/go.mod src/go.sum ./
RUN go mod download
COPY src/ .
RUN go build -o /api main.go
CMD ["/api"]
