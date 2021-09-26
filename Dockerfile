FROM golang:latest

WORKDIR /go/src/app

#COPY ./ ./
ADD . .
#RUN go mod init

RUN go mod download
RUN go build -o main .

#EXPOSE 80:8080
EXPOSE 6111

CMD ["./main"]