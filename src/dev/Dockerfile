FROM golang:1.17

ARG DEBIAN_FRONTEND=noninteractive
ENV TZ=Asia/Singapore

RUN apt-get update -y && \
apt-get upgrade -y && \
apt-get install -y \
tzdata \
sudo \
vim

RUN useradd dev -ms /bin/bash -G sudo
RUN echo "\ndev\tALL=(ALL) NOPASSWD:ALL\n" >> /etc/sudoers
WORKDIR /apps
RUN mkdir echo && mkdir generator && chown -R dev:dev .

USER dev

WORKDIR /home/dev
RUN git clone https://github.com/fatih/vim-go.git .vim/pack/plugins/start/vim-go
RUN vim -esN +GoInstallBinaries +q
RUN vim -esN +"helptags ALL" +q
COPY .vimrc .

WORKDIR /apps/echo
RUN go mod init echo
RUN go get github.com/labstack/echo/v4 
RUN go get github.com/go-sql-driver/mysql

WORKDIR /apps/generator
RUN go mod init generator
RUN go get github.com/labstack/echo/v4 

WORKDIR /apps

CMD ["/bin/bash"]
