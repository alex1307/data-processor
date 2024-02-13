# Use the official Go image as a parent image.
FROM golang:latest as builder


# Set the working directory inside the container.
WORKDIR /app

# Copy the go.mod and go.sum files.
COPY go.mod go.sum ./
# Download dependencies.
RUN go mod download

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
FROM scratch
COPY --from=builder /app/processor /bin/processor


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
