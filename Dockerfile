FROM golang:1.21-alpine3.18 as builder
LABEL maintainer="manishkp220@gmail.com"

RUN mkdir -p /opt/pod-dependency
WORKDIR /opt/pod-dependency

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .


######## Start a new stage from scratch #######
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /opt/pod-dependency

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /opt/pod-dependency  .

# Expose port 8080 to the outside world
EXPOSE 8090

# Command to run the executable
CMD ["./main"]