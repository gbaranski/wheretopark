FROM golang:1.19.1-alpine as BUILD
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY authorization/**.go ./
RUN go build -o /bin/my-server .

FROM alpine:latest
RUN apk --no-cache add curl
COPY --from=BUILD /bin/my-server /bin/my-server
CMD [ "/bin/my-server" ]

