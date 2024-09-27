FROM golang:1.22.6-alpine3.20 AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o resolver ./cmd/app/main.go

FROM alpine:3.20

RUN addgroup -S appgroup && adduser -S appuser -G appgroup
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/resolver /usr/local/bin/resolver

RUN chown appuser:appgroup /usr/local/bin/resolver

WORKDIR /usr/local/bin

USER appuser

CMD ["resolver"]