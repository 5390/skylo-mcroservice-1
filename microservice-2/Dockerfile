# Base image
FROM golang:1.20-alpine

# Set the working directory
WORKDIR /app

# Copy the source code
COPY . .

# Install dependencies
RUN go mod tidy

# Build the application
RUN go build -o microservice-2 main.go

# Expose the service port
EXPOSE 8081

# Run the service
CMD ["./microservice-2"]
