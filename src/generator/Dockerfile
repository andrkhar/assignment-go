FROM golang:1.17

ARG DEBIAN_FRONTEND=noninteractive
ENV TZ=Asia/Singapore

RUN apt-get update -y && \
apt-get upgrade -y && \
apt-get install -y \
tzdata

WORKDIR /apps/generator
RUN useradd -m generator
RUN chown -R generator:generator .

USER generator
RUN go mod init echo
RUN go get github.com/labstack/echo/v4 

COPY . .

CMD ["go", "run", "server.go"]
