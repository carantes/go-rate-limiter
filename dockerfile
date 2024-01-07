FROM golang:1.21-alpine

RUN mkdir /goratelimiter
WORKDIR /goratelimiter

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o /go/bin/goratelimiter

ENTRYPOINT [ "/go/bin/goratelimiter" ]
EXPOSE 8080