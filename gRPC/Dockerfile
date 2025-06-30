FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o server .

FROM alpine
WORKDIR /root/
COPY --from=builder /app/server .

CMD ["./server"]
