FROM golang:1.22

WORKDIR /app

COPY . .
RUN go mod download && go mod verify

RUN CGO_ENABLED=0 go build -o bin/app

ENV PORT 5000
ENV GIN_MODE release
EXPOSE 5000

ENTRYPOINT ["go", "run", "main.go"]