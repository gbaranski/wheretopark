FROM golang:1.19.1-alpine
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY authorization/**.go ./
RUN go build -o /bin/authorization-server .
CMD [ "/bin/authorization-server" ]

