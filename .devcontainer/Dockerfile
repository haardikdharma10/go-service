FROM golang:1.22.4

RUN apt-get update && apt-get install -y sqlite3 libsqlite3-dev

WORKDIR /workspace

COPY ../go.mod ../go.sum ./

RUN go mod download

COPY .. .

EXPOSE 8080