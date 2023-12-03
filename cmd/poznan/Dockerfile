FROM golang:1-alpine as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY go ./go
COPY ./collector/poznan ./collector/poznan
COPY ./cmd/poznan ./cmd/poznan

RUN go build -v -o exe ./cmd/poznan

FROM golang:1-alpine

COPY --from=builder /app/exe /app/exe
CMD ["/app/exe"]