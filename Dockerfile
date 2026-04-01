FROM golang:1.25-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/goadmin .

FROM alpine:3.22

WORKDIR /app

RUN apk add --no-cache tzdata \
 && adduser -D -u 10001 appuser

COPY --from=builder /out/goadmin /app/goadmin
COPY configs /app/configs
COPY docs /app/docs

USER appuser

EXPOSE 3000

ENTRYPOINT ["/app/goadmin"]
CMD ["start", "-d", "/app/configs", "-c", "docker"]
