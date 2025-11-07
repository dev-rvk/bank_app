# Build stage - install go deps and generate an executable 
FROM golang:1.24.5-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go


# Run Stage - just execute the binary (no need to have dependencies)
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .

EXPOSE 8080
CMD [ "/app/main" ]