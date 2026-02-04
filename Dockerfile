FROM golang:1.25.6-alpine AS builder

RUN apk --no-cache add git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
RUN go mod tidy

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o bin cmd/api/main.go

FROM alpine AS runner

WORKDIR /app

COPY --from=builder /app/bin .

# Открываем порт
EXPOSE ${HTTP_PORT}

CMD [ "/app/bin" ]