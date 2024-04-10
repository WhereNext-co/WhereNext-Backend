# Build stage
FROM golang:1.22 AS build-env

WORKDIR /app

COPY . .
RUN go mod download && go mod verify
RUN CGO_ENABLED=0 go build -o bin/app

# Production stage
FROM alpine:latest

WORKDIR /app

COPY --from=build-env /app/bin/app /app/app
COPY --from=build-env /app/.env /app/.env

ENV PORT 5000
ENV GIN_MODE release
EXPOSE 5000

ENTRYPOINT ["/app/app"]