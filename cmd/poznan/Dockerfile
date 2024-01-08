FROM golang:1-alpine as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY go ./go
COPY ./providers/poznan ./providers/poznan
COPY ./cmd/poznan ./cmd/poznan

RUN go build -v -o exe ./cmd/poznan

FROM golang:1-alpine

COPY --from=builder /app/exe /app/exe
CMD ["/app/exe"]