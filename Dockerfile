FROM golang:1.22-bullseye

EXPOSE 25
EXPOSE 8080

ENV SLACK_URL=""
ENV DOMAIN=""
ENV USERNAME=""
ENV SECRET=""

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o app

CMD [ "sh", "-c", "./app -s $SLACK_URL -d $DOMAIN -t $SECRET" ]
