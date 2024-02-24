FROM golang:1.18-bullseye

ENV ADDR="0.0.0.0:25" \
    SLACK_URL="" \
    DOMAIN="" \
    USERNAME=""

WORKDIR /app
COPY . .

RUN go mod download

RUN go build cmd/serve.go

EXPOSE 25

CMD [ "/serve", "-l", ${ADDR}, "-s", ${SLACK_URL}, "-d", ${DOMAIN}, "-u", ${USERNAME} ]
