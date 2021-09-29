#FROM golang:latest
#
#WORKDIR /go/src/app
#
##COPY ./ ./
#ADD . .
##RUN go mod init
#
#RUN go mod download
#RUN go build -o main .
#
##EXPOSE 80:8080
#EXPOSE 6111
#
#CMD ["./main"]


#FROM golang:latest
#
#RUN go version
#ENV GOPATH=/
#
#COPY ./ ./
#
#RUN go mod download
#RUN go built -o CatsCrud ./main.go
#
#CMD ["go run", "main.go"]

FROM golang:latest

RUN mkdir -p usr/src/app/
WORKDIR /usr/src/app/

COPY . /usr/src/app/

# install psql
RUN apt-get update
RUN apt-get -y install postgresql-client

RUN chmod +x wait-for-postgres.sh

RUN go mod download
RUN go build -o app main.go

CMD ["./app"]