FROM gocv/opencv:4.6.0 as builder

ENV GOPATH /go

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./
RUN go version

RUN go build -v -o cctv ./cmd/cctv

FROM gocv/opencv:4.6.0

COPY --from=builder /app/cctv /app/cctv
# ENTRYPOINT ldd /app/cctv 
CMD ["/app/cctv"]