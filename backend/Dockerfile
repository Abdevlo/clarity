FROM golang:1.21-alpine AS builder

WORKDIR /build

RUN apk add --no-cache protobuf-dev

COPY . .

RUN go mod download

RUN make proto

RUN CGO_ENABLED=1 GOOS=linux go build -o clarity-backend main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /build/clarity-backend .
COPY --from=builder /build/.env .env

EXPOSE 50051

CMD ["./clarity-backend"]
