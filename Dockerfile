# BUILDING THE APP
FROM golang:1.21.4-alpine AS builder

# set the current Working Directory inside the container
RUN mkdir /app
WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# This will copy all the files in our repo to the inside the container at root location.
COPY . .
# Build our binary at ./cmd/main.go location.
RUN go build -o /tigerhall-kittens ./cmd/main.go

# DEPLOYING
FROM alpine:latest

# copy the already-built binary from the builder, then run it
WORKDIR /
COPY --from=builder /tigerhall-kittens /tigerhall-kittens
EXPOSE 8080
ENTRYPOINT ["/tigerhall-kittens"]


# Set the local PORT environment variable inside the container
ENV PORT 8080
ENV DB_HOST postgres
ENV DB_USER postgres
ENV DB_PASS password
ENV DB_NAME tigerhall_kittens
ENV DB_PORT 5432
ENV MAILTRAP_USERNAME 8d966b48a9e1fa
ENV MAILTRAP_PASSWORD 231130ffcbf49f
ENV SECRET jwt_secret
