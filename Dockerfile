FROM golang as build

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /bin/bot

FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /bin/bot /bin/bot
COPY --from=build /app/wait-for-it.sh /bin/wait-for-it
ENTRYPOINT ["./wait-for-it.sh","es01:9200","--","/bin/bot"]