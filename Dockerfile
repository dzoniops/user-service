# Build the application from source
FROM golang:1.19 AS BuildStage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /user-service main.go

# Run the tests in the container
FROM alpine:latest

WORKDIR /

COPY --from=BuildStage /user-service /user-service

EXPOSE 8080 

USER nonroot:nonroot

ENTRYPOINT ["/user-service"]