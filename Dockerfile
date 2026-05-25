# ===== BUILD STAGE =====
FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY . .

RUN go mod tidy

# Compila o executável real baseado no main.go da raiz
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o security-dashboard .

# ===== RUNTIME STAGE =====
FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates nmap nmap-scripts

# Copia o binário com a permissão correta
COPY --chmod=755 --from=builder /app/security-dashboard .

# Assets
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static

EXPOSE 8080

CMD ["./security-dashboard"]
