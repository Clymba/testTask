FROM golang:1.21.2-alpine3.18 AS builder
WORKDIR /app
COPY tools /app
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/app.go

FROM alpine:3.18
COPY --from=builder /app/main /app/service
ENTRYPOINT [ "/app/service" ]