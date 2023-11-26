# Dockerfile definition for Backend application service.

####################################################################
# We're adopting multi-stage Dockerfile. This is the first stage
# to build the application binary. This is basically our environment.
FROM golang:1.19-alpine as build

WORKDIR /usr/src/app

# Precache dependencies.
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Build the app binary.
COPY . .
RUN go build -v -o /usr/bin/app ./cmd/.

####################################################################
# This is the actual image that we will be using in production.
FROM alpine:latest

# Copy app binary and certificate files from builder stage.
COPY --from=build /usr/bin/app /usr/bin/
COPY --from=build /usr/src/app/*.pem /etc/app/

WORKDIR /root

# This is the port that our application will be listening on.
EXPOSE 1323

# Execute app binary when the container is started.
ENTRYPOINT ["/usr/bin/app"]