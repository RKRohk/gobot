FROM golang

WORKDIR /go/src/github.com/rkrohk/gobot/

COPY . .


CMD [ "while true; do foo; sleep 2; done" ]