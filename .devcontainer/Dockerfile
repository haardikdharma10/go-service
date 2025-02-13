FROM golang:1.22.4

# Install SQLite
RUN apt-get update && apt-get install -y sqlite3 libsqlite3-dev

WORKDIR /workspace

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the code
COPY . .

# Expose the port the app runs on
EXPOSE 8080