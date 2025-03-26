FROM golang:1.23-alpine
WORKDIR /app2
COPY . .
RUN go mod download || echo "No go.mod, skipping dependencies"
RUN go build -v -o counter || { echo "Build failed"; exit 1; }
EXPOSE 8081
CMD ["./counter"]