FROM golang:1.17-alpine AS builder

LABEL maintainer="barikly <barikly@pu.go.id>"

# Move to working directory (/build).
WORKDIR /build

# Copy and download dependency using go mod.
COPY go.mod go.sum .env ./
RUN go mod download

# Copy the code into the container.
COPY . .

# Set necessary environment variables needed for our image and build the API server.
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -o apiserver .


# Copy binary and config files from /build to root folder of scratch container.
FROM gcr.io/distroless/static-debian11
COPY --from=builder ["/build/apiserver", "/build/.env", "/"]
CMD ["./apiserver"]
