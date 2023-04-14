FROM golang:1.20-buster

MAINTAINER daniel pickens

WORKDIR /go


RUN echo $PATH
RUN ls
RUN pwd

CMD ["./hercules"]
