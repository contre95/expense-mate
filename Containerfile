FROM golang:1.20
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY cmd pkg .
RUN go build -o /app/app ./cmd/api/main.go
CMD ["/app/app"]
