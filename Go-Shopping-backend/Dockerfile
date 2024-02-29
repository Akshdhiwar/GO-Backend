# Start from the official Golang image
FROM golang:latest AS build-stage

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

COPY views /app/views

# Build the Go app
RUN go build -o main .

FROM golang:latest AS release-stage

WORKDIR /

COPY --from=build-stage /app/main /main

# Copy the views directory from the build stage
COPY --from=build-stage /app/views /views

# Expose port 3000 to the outside world
EXPOSE 3000

# Command to run the executable
ENTRYPOINT ["/main"]
