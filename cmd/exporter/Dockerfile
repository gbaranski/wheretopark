FROM golang:1-alpine as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY ./cmd/exporter ./cmd/exporter
COPY go ./go
COPY ./prometheus/exporter ./prometheus/exporter

RUN go build -v -o exporter ./cmd/exporter

FROM golang:1-alpine

COPY --from=builder /app/exporter /app/exporter
CMD ["/app/exporter"]