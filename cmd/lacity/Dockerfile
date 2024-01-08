FROM golang:1-alpine as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY go ./go
COPY ./providers/lacity ./providers/lacity
COPY ./cmd/lacity ./cmd/lacity

RUN go build -v -o exe ./cmd/lacity

FROM golang:1-alpine

COPY --from=builder /app/exe /app/exe
CMD ["/app/exe"]