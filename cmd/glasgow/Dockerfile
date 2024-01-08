FROM golang:1-alpine as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY go ./go
COPY ./providers/glasgow ./providers/glasgow
COPY ./cmd/glasgow ./cmd/glasgow

RUN go build -v -o exe ./cmd/glasgow

FROM golang:1-alpine

COPY --from=builder /app/exe /app/exe
CMD ["/app/exe"]