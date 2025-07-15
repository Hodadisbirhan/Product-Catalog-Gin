FROM golang:1.24-alpine


# Install required packages
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod .
COPY go.sum .

COPY vendor/ ./vendor/

ENV GOPROXY=direct
ENV GOSUMDB=off
# Download dependencies
# RUN go mod download

# Copy source code
COPY . .

# Build the Go app
RUN go build -mod=vendor -o main .

# Expose port 8080
EXPOSE 8080

# Command to run the executable
CMD ["./main"]