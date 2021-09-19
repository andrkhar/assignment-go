FROM golang:1.17

ARG DEBIAN_FRONTEND=noninteractive
ENV TZ=Asia/Singapore

RUN apt-get update -y && \
apt-get upgrade -y && \
apt-get install -y \
tzdata

WORKDIR /apps/echo
RUN useradd -m echo
RUN chown -R echo:echo .

USER echo
RUN go mod init echo
RUN go get github.com/labstack/echo/v4 
RUN go get github.com/go-sql-driver/mysql

COPY . .

CMD ["go", "run", "server.go"]
