FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download || true
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/app ./cmd

FROM alpine:latest
WORKDIR /app
COPY --from=builder /out/app /app/app
# Copia arquivos auxiliares se existirem (ex: zones.json)
COPY . /app/
CMD ["/app/app"]
