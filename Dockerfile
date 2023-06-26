FROM golang:1.19-bullseye

WORKDIR /app

COPY . .

RUN go build -o linxdatacenter main.go

ENTRYPOINT ["./linxdatacenter"]