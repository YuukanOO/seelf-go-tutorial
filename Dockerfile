FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN go build -ldflags="-s -w" -o sample

FROM alpine:3.16
WORKDIR /app
COPY --from=builder /app/sample ./
EXPOSE 8080
CMD ["./sample"]