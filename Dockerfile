FROM golang:1.21.3-alpine3.18

COPY . /app
WORKDIR /app

RUN mkdir /app/bin
RUN mkdir /app/public

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/🍍

EXPOSE 8000

CMD /app/bin/🍍 -port 8000 -public /app/public
