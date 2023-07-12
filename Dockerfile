FROM golang:1.19-alpine as build

ADD . /build
WORKDIR /build

RUN go build -o /build/webhook -ldflags "-s -w"


FROM alpine:3.16

COPY --from=build /build/webhook /srv/webhook
COPY ./config.yaml /srv/
RUN chmod +x /srv/webhook

WORKDIR /srv
EXPOSE 8080
CMD ["/srv/webhook"]