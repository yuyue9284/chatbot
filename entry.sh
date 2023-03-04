#!/bin/sh
openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout /app/key.pem -out /app/cert.pem -subj "/C=US/ST=Denial/L=Springfield/O=Dis/CN=testyy.com"
./gotty -w --tls --tls-crt /app/cert.pem --tls-key /app/key.pem /app/main
