FROM golang:1-alpine as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY ./cmd/collector ./cmd/collector
COPY go ./go
COPY collector ./collector

RUN go build -v -o collector ./cmd/collector

FROM golang:1-alpine

COPY --from=builder /app/collector /app/collector
CMD ["/app/collector"]