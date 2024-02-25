FROM golang:1.19-bullseye

EXPOSE 25

ENV SLACK_URL=""
ENV DOMAIN=""
ENV USERNAME=""

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o app cmd/serve.go

CMD [ "sh", "-c", "./app -l 0.0.0.0:25 -s $SLACK_URL -d $DOMAIN" ]
