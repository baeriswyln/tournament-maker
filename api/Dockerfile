# Sources:
#   - https://docs.docker.com/language/golang/build-images/
#   - https://snyk.io/blog/containerizing-go-applications-with-docker/
# syntax=docker/dockerfile:1

##
## STEP 1 - BUILD
##
FROM golang:1.26-alpine AS builder

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY . ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /build ./cmd/server



##
## STEP 2 - DEPLOY
##
FROM alpine:3.19

# Set destination for COPY
WORKDIR /

# Copy binaries
COPY --from=builder /build /build

# Expose the API port
EXPOSE 8080

# Run
CMD ["/build"]