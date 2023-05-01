FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./social ./cmd/app/main.go

FROM scratch
COPY --from=builder /app/social /usr/bin/social
COPY --from=builder /app/internal/migrations /usr/bin/migrations
ENTRYPOINT [ "/usr/bin/social" ]