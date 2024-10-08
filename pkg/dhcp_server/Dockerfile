# Cria a imagem builder para compilar o código fonte em um executável
FROM golang:1.22.2-alpine as builder
RUN apk add --no-cache gcc musl-dev
WORKDIR /app

# Copia os arquivos de dependências para o diretório de trabalho
COPY go.mod go.sum ./
RUN go mod download

# Copia todo o código do projeto para o diretório de trabalho
COPY . .

# Compila o aplicativo Go
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o dhcp_server ./cmd/main.go

# Cria a imagem final, executando a API e expondo a porta 8080
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/dhcp_server .
ARG PORT=67
ENV PORT=$PORT
EXPOSE $PORT
CMD ["./dhcp_server"]
