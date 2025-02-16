# Stage 1: Build the Go binary
FROM golang:alpine AS builder

# Required for alpine + sqlite3 driver
ENV CGO_ENABLED=1
RUN apk add --no-cache gcc musl-dev

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the source from the current directory to the Working Directory inside the container
ADD . .

# Tidy up and download dependencies
RUN go mod tidy

# Build the Go app with necessary flags for static linking
RUN go build -ldflags='-s -w -extldflags "-static"' -o /app/expenses-app cmd/main.go

# Stage 2: Run the Go binary in a minimal environment
FROM scratch

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the required files and directories from the build stage
COPY ./views /app/views
COPY ./public /app/public
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/expenses-app /app/expenses-app

# Expose the port the app runs on
EXPOSE 3535

# Command to run the executable
ENTRYPOINT ["/app/expenses-app"]

