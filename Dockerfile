FROM golang:1.24-alpine

# Install required packages
RUN apk add --no-cache git curl

# Set environment variables
ENV GOPROXY=direct
ENV GOSUMDB=off
ENV CGO_ENABLED=0

# Install Air
RUN go install github.com/air-verse/air@latest


# Set working directory
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod ./
COPY go.sum ./
COPY vendor/ ./vendor/

# Copy the source code
COPY . .

# Expose the port the app runs on
EXPOSE 8080

# Run Air for hot reload
CMD ["air"]
