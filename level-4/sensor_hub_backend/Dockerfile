FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /api ./main.go

# Stage 2: Runtime
FROM alpine:3.20

WORKDIR /root/

COPY --from=builder /api .

EXPOSE 8080

CMD ["./api"]