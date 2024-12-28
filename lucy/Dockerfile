FROM golang:1.23-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

RUN sqlc generate

RUN go build -o lucy ./cmd

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/lucy .

EXPOSE 3333

CMD ["./lucy"]
