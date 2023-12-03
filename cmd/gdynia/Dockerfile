FROM golang:1-alpine as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY go ./go
COPY ./collector/gdynia ./collector/gdynia
COPY ./cmd/gdynia ./cmd/gdynia

RUN go build -v -o exe ./cmd/gdynia

FROM golang:1-alpine

COPY --from=builder /app/exe /app/exe
CMD ["/app/exe"]