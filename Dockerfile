# Use the official Go 1.22 image as the base image
FROM golang:1.22-alpine

# Install required dependencies for CGO and SQLite
RUN apk add --no-cache gcc musl-dev sqlite-dev

# Enable CGO
ENV CGO_ENABLED=1

# Set the working directory
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application
COPY . .

# Build the Go application, specifying the correct main file (bot.go)
RUN go build -o app bot.go

# Command to run the Go application
CMD ["./app"]