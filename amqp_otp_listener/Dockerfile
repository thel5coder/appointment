FROM golang:1.14.0

RUN apt-get update && apt-get install -y

ENV GO111MODULE=on
ENV PKG_NAME=profira-backend/
ENV PKG_PATH=$GOPATH/src/$PKG_NAME
WORKDIR $PKG_PATH/

COPY . $PKG_PATH/

RUN echo $PWD
RUN go mod vendor

WORKDIR $PKG_PATH/amqp_otp_listener/
RUN echo $PWD

RUN go build main.go
EXPOSE 3001
CMD ["./main"]