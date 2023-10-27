FROM golang:1.21.3-alpine3.18 as build

COPY . /app
WORKDIR /app

RUN mkdir /app/bin
RUN mkdir /app/public

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/🍍

FROM gcr.io/distroless/static-debian11

COPY --from=build /app/bin/🍍 /bin

EXPOSE 8000

CMD ["/bin/🍍", "-port", "8000", "-public", "/app/public"]
