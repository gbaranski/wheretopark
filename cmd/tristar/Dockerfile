FROM golang:1-alpine as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./

RUN go build -v -o tristar ./cmd/tristar

FROM golang:1-alpine

COPY --from=builder /app/tristar /app/tristar
CMD ["/app/tristar"]