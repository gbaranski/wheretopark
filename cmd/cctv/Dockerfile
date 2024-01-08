FROM gocv/opencv:4.8.0 as builder

ENV GOPATH /go

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY go ./go
COPY ./providers/cctv ./providers/cctv
COPY ./cmd/cctv ./cmd/cctv

RUN go build -v -o exe ./cmd/cctv

FROM gocv/opencv:4.8.0

RUN apt-get update && apt-get -y install ffmpeg

COPY --from=builder /app/exe /app/exe

CMD ["/app/exe"]
