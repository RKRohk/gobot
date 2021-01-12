FROM golang

WORKDIR /app

COPY go.* ./

COPY . .

RUN go build


CMD go run bot.go