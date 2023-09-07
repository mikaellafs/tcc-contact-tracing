FROM golang:1.18

WORKDIR /app

COPY . .

RUN go get -d -v ./pkg/cmd
RUN go build -o server -v ./pkg/cmd

CMD ["./server"]