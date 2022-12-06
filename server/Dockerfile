FROM golang:1.18

WORKDIR /app

COPY . .

RUN go get -d -v ./src/cmd
RUN go build -o server -v ./src/cmd

CMD ["./server"]