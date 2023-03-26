FROM golang:1-alpine as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY ./cmd/influxer ./cmd/influxer
COPY go ./go

RUN go build -v -o influxer ./cmd/influxer

FROM golang:1-alpine

COPY --from=builder /app/influxer /app/influxer
CMD ["/app/influxer"]