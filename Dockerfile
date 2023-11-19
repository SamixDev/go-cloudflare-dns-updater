# syntax=docker/dockerfile:1

FROM golang:1.21.1 as Base

# Set destination for COPY
WORKDIR /

# Download Go modules
COPY go.mod go.sum .env ./
RUN go mod download
# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY *.go ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-gs-ping

# Run
CMD ["/docker-gs-ping"]