FROM golang:1-alpine as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY go ./go
COPY ./cmd/influxer ./cmd/influxer


RUN go build -v -o exe ./cmd/influxer

FROM golang:1-alpine

COPY --from=builder /app/exe /app/exe
CMD ["/app/exe"]