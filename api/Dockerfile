FROM golang:latest
LABEL maintainer="github.com/innovember"

WORKDIR /go/src/forum

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

EXPOSE 8081

RUN go build -o main .

CMD ["./main"]
