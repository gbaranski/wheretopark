FROM golang:1-alpine as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY go ./go
COPY ./collector/gdansk ./collector/gdansk
COPY ./cmd/gdansk ./cmd/gdansk

RUN go build -v -o exe ./cmd/gdansk

FROM golang:1-alpine

COPY --from=builder /app/exe /app/exe
CMD ["/app/exe"]