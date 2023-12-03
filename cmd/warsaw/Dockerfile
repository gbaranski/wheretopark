FROM golang:1-alpine as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY go ./go
COPY ./collector/warsaw ./collector/warsaw
COPY ./cmd/warsaw ./cmd/warsaw

RUN go build -v -o exe ./cmd/warsaw

FROM golang:1-alpine

COPY --from=builder /app/exe /app/exe
CMD ["/app/exe"]