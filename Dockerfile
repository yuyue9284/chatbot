FROM golang:1.17 AS builder
WORKDIR /app
COPY . .
RUN go build -o /app/main .
RUN wget https://github.com/yudai/gotty/releases/download/v2.0.0-alpha.3/gotty_2.0.0-alpha.3_linux_amd64.tar.gz -qO- | tar -xz

FROM debian AS runner
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/gotty .
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# create self-signed certificate for gotty
RUN openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout /app/key.pem -out /app/cert.pem -subj "/C=US/ST=Denial/L=Springfield/O=Dis/CN=testyy.com"
CMD ["./gotty", "-w", "--tls", "--tls-crt", "/app/cert.pem", "--tls-key", "/app/key.pem", "./main"]
