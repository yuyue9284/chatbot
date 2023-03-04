FROM golang:1.17 AS builder
WORKDIR /app
COPY . .
RUN go build -o /app/main .
RUN wget https://github.com/yudai/gotty/releases/download/v2.0.0-alpha.3/gotty_2.0.0-alpha.3_linux_amd64.tar.gz -qO- | tar -xz

FROM debian AS runner
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/gotty .
COPY ./entry.sh .
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/* && chmod +x /app/entry.sh
EXPOSE 8080
ENTRYPOINT ["/app/entry.sh"]