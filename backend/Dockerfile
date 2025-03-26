FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o backend

FROM alpine:latest
RUN apk add --no-cache \
  chromium \
  nss \
  freetype \
  harfbuzz \
  ttf-freefont \
  && rm -rf /var/cache/apk/*
COPY --from=builder /app/backend /usr/local/bin/backend
EXPOSE 8080
CMD ["backend"]