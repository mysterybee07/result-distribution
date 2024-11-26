# Use the official Go image with Alpine variant for a lightweight container
FROM golang:1.23-alpine

# Set the working directory inside the container to /app
WORKDIR /app

# Install air for live reloading
RUN go install github.com/air-verse/air@latest

# Copy the .air.toml file into the container
COPY .air.toml /app/

# Copy go.mod and go.sum for dependency management
COPY go.mod go.sum ./

# Download the Go module dependencies specified in go.mod
RUN go mod download

# Copy all source files and subdirectories into the working directory
COPY . .

# Ensure proper permissions for all files (if needed)
RUN chmod -R 755 /app

# Expose port 8080 for the application to listen on
EXPOSE 8080

# Run air with configuration when the container starts
CMD ["air", "-c", ".air.toml"]
