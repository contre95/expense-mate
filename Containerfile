# Stage 1: Build the Go binary
FROM golang:alpine AS builder

# Required for alpine + sqlite3 driver
ENV CGO_ENABLED=1
RUN apk add --no-cache gcc musl-dev

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy only the files needed for dependency resolution
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go app with necessary flags for static linking
RUN go build -trimpath -ldflags='-s -w -extldflags "-static"' -o /app/expenses-app cmd/main.go

# Stage 2: Run the Go binary in a minimal environment
FROM scratch
# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the required files and directories from the build stage
COPY --from=builder /app/expenses-app /app/expenses-app
COPY --from=builder /app/views /app/views
COPY --from=builder /app/public /app/public
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Expose the port the app runs on
EXPOSE 8080

# Command to run the executable
ENTRYPOINT ["/app/expenses-app"]
