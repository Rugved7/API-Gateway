FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY go.mod go.sum ./
RUN go mod download



COPY . .


RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o gateway ./cmd/gateway



FROM gcr.io/distroless/base-debian12

WORKDIR /app


COPY --from=builder /app/gateway /gateway


COPY configs /configs


EXPOSE 8080


ENV GATEWAY_CONFIG=/configs/gateway.yaml


USER nonroot:nonroot

ENTRYPOINT ["/gateway"]
