FROM golang:1-alpine as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY go ./go
COPY ./meters ./meters
COPY ./collector/krakow ./collector/krakow
COPY ./forecaster ./forecaster
COPY ./cmd/forecaster ./cmd/forecaster

RUN go build -v -o exe ./cmd/forecaster

FROM golang:1-alpine

COPY --from=builder /app/exe /app/exe
CMD ["/app/exe"]