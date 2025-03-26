FROM golang:1.23-alpine
WORKDIR /app1
COPY . .
RUN go mod download || echo "No go.mod, skipping dependencies"
RUN go build -v -o bot || { echo "Build failed"; exit 1; }
CMD ["./bot"]