# Use the official Golang image to build the Go application
FROM golang:1.24.2 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application
RUN make build

# Use a minimal image for the final container
FROM node:20

# Set the working directory inside the container
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/backend-exec .

# Expose the port the app runs on
EXPOSE 3500

# install pm2
RUN npm install pm2 -g

# create empty .env file to avoid error
RUN touch .env

CMD ["pm2-runtime", "start", "backend-exec"]
