FROM golang:1.26-alpine AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/app ./cmd/main.go

FROM gcr.io/distroless/static:nonroot
WORKDIR /app
COPY --from=build /out/app ./main
USER nonroot:nonroot
ENTRYPOINT ["/app/main"]
