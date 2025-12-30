FROM golang:1.25.2-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	go build -o app

FROM alpine:latest AS runtime

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/app .
COPY --from=builder /app/migrations ./migrations


EXPOSE 8080

CMD ["./app"]
	