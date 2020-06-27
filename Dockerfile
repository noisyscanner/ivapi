FROM golang:1.14.4-alpine3.12 AS dev

RUN apk add --no-cache dep git make
RUN go get -u github.com/cosmtrek/air

ARG ROOT=/go/src/bradreed.co.uk/iverbs/api
WORKDIR $ROOT

COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure -vendor-only

COPY . .

ENV REDIS=host.docker.internal:6379 \
  PORT=7000

EXPOSE $PORT

CMD ["make", "run"]

FROM dev AS build

RUN make build

FROM alpine:3.12 AS prod

WORKDIR /usr/local/bin

ARG ROOT=/go/src/bradreed.co.uk/iverbs/api
COPY --from=build $ROOT/tmp/api .

CMD api
