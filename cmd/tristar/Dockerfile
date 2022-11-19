FROM golang:1-alpine as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY ./cmd/tristar ./cmd/tristar
COPY go ./go
COPY providers/tristar ./providers/tristar

RUN go build -v -o tristar ./cmd/tristar

FROM golang:1-alpine

COPY --from=builder /app/tristar /app/tristar
CMD ["/app/tristar"]