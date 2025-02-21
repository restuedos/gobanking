FROM golang:1.24-alpine

WORKDIR /app

# Install air
RUN go install github.com/air-verse/air@latest

# Install swag
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Generate swagger documentation
RUN swag init

# Build the application
RUN go build -o main .

# Expose port
EXPOSE 3000

# Command to run the executable
CMD ["air", "-c", ".air.toml"]