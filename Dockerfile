FROM golang:1.13-alpine3.10 as build

RUN apk add --no-cache --update git make

RUN mkdir /build
WORKDIR /build
COPY . .
RUN make build


FROM alpine:3.10

ARG REVISION
ARG VERSION

LABEL maintainer="Rico Berger"
LABEL git.url="https://github.com/ricoberger/logmetrics"

RUN apk add --no-cache --update curl ca-certificates

USER nobody

COPY --from=build /build/bin/logmetrics /bin/logmetrics
EXPOSE 8080

ENTRYPOINT  [ "/bin/logmetrics" ]
