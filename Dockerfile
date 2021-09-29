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