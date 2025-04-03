FROM golang:1.24-alpine

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o main .

# Create .env file from example if it doesn't exist
RUN cp -n .env.example .env 2>/dev/null || true

# Expose port 8080
EXPOSE 8080

# Command to run the application
CMD ["./main"]
