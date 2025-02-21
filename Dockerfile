FROM golang:latest as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o chi-crud-api main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/chi-crud-api .

EXPOSE 8080

CMD ["./chid-crud-api"]