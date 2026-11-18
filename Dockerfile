# syntax=docker/dockerfile:1

FROM golang:1.26-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# TARGETOS/TARGETARCH are set by BuildKit, so the image builds natively on
# both amd64 and arm64 hosts.
ARG TARGETOS=linux
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
	go build -trimpath -ldflags="-s -w" -o /out/pulse-api .

FROM alpine:3.22 AS runtime

RUN apk add --no-cache ca-certificates tzdata \
	&& adduser -D -H -u 10001 pulse

WORKDIR /app

# Migrations are embedded in the binary.
COPY --from=builder /out/pulse-api /app/pulse-api

USER pulse

EXPOSE 8080

ENTRYPOINT ["/app/pulse-api"]
