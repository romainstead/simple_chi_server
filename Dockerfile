FROM golang:latest as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o chi-crud-api .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/chi-crud-api /app/chi-crud-api

EXPOSE 8080

CMD ["./chi-crud-api"]