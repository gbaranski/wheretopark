FROM golang:1-bullseye as builder

RUN apt-get update &&  apt-get install -y \
    pkg-config \
    libopencv-dev

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./

RUN go build -v -o cctv ./cmd/cctv

FROM golang:1-bullseye

RUN apt-get update &&  apt-get install -y \
    libopencv-dev

COPY --from=builder /app/cctv /app/cctv
CMD ["/app/cctv"]