FROM golang as build

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /bin/bot

FROM scratch
COPY --from=build /bin/bot /bin/bot
ENTRYPOINT ["/bin/bot"]