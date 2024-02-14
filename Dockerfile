# Use the official Go image as a parent image.
FROM golang:1.22.0-alpine3.19 as builder
RUN apk add --no-cache musl-dev
RUN apk add --no-cache openssl-dev
RUN apk add --no-cache pkgconfig
RUN apk update && apk add bash
RUN apk add --no-cache gcc
RUN apk --no-cache add \
    bash \
    g++ \
    ca-certificates \
    lz4-dev \
    musl-dev \
    cyrus-sasl-dev \
    openssl-dev \
    make \
    python3 \
    protoc 


# Set the working directory inside the container.
WORKDIR /app

# Copy the go.mod and go.sum files.
COPY go.mod go.sum ./
# Download dependencies.
RUN go mod download
RUN go mod tidy
# Copy the rest of the source code.
COPY . ./
COPY resources/config /app/resources/config

# Use the Makefile to build the application.
RUN echo "Current directory:" && pwd
RUN echo "Contents of current directory:" && ls -la
RUN go clean
RUN go build -o processor ./cmd/processor/main.go
RUN chmod +x ./processor

# Deploy the application binary into a lean image
FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN mkdir /bin/logs
COPY --from=builder /app/processor /bin/processor
COPY --from=builder /app/resources/config /bin/resources/config

ENV KAFKA_BROKER=kafka:9092
ENV DB_HOST=postgres-server
ENV DB_USER=admin
ENV DB_PASSWORD=1234
ENV DB_PORT=5432
ENV SSL_MODE=disable
ENV DB_SCHEMA=public
ENV DB_NAME=vehicles



# Command to run the compiled binary.
CMD ["bin/processor"] 
