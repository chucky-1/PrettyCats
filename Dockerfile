FROM golang:latest

#ADD docker-entrypoint-initdb.d/init.sh /usr/local/bin/docker-entrypoint.sh/docker-entrypoint-initdb.d/
#COPY docker-entrypoint-initdb.d /docker-entrypoint-initdb.d

RUN mkdir -p usr/src/app/
WORKDIR /usr/src/app/

COPY . /usr/src/app/

RUN go mod download

RUN go build -o /docker-cats

CMD ["/docker-cats"]