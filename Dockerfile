FROM golang:1.21-alpine

WORKDIR /app

# Install swag
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy go mod and sum files
COPY go.mod ./

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
CMD ["./main"]