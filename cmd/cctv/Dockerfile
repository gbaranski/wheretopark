FROM gocv/opencv:4.8.0 as builder

ENV GOPATH /go

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY ./cmd/cctv ./cmd/cctv
COPY go ./go
COPY providers/cctv ./providers/cctv

RUN go build -v -o cctv ./cmd/cctv

FROM gocv/opencv:4.8.0

RUN apt-get update && apt-get -y install ffmpeg

COPY --from=builder /app/cctv /app/cctv

CMD ["/app/cctv"]