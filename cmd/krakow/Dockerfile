FROM golang:1-alpine as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY go ./go
COPY ./collector/krakow ./collector/krakow
COPY ./cmd/krakow ./cmd/krakow

RUN go build -v -o exe ./cmd/krakow

FROM golang:1-alpine

COPY --from=builder /app/exe /app/exe
CMD ["/app/exe"]
