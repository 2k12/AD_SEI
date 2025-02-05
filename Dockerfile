ARG GO_VERSION=1
FROM golang:${GO_VERSION}-bookworm AS builder

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download
COPY . .  
RUN go build -o run-app . 

# Imagen final
FROM debian:bookworm

WORKDIR /app
COPY --from=builder /usr/src/app/run-app /app/  
COPY --from=builder /usr/src/app/docs /app/docs 

COPY --from=builder /usr/src/app/assets/fonts /app/assets/fonts

ENV PORT=8080

EXPOSE 8080

CMD ["./run-app"]
