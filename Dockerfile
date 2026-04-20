FROM golang:1.26-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o attengo ./cmd/attengo

FROM alpine:3.20

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/attengo .

COPY --from=builder /app/db/migrations ./db/migrations

COPY --from=builder /app/web/templates ./web/templates
COPY --from=builder /app/web/static ./web/static

RUN mkdir -p /data

EXPOSE 8080

CMD ["./attengo"]