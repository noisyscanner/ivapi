FROM golang:1.14.4-alpine3.12 AS dev

RUN apk add --no-cache git make
RUN go get -u github.com/cosmtrek/air

ARG ROOT=/go/src/github.com/noisyscanner/ivapi
WORKDIR $ROOT

ENV REDIS=host.docker.internal:6379 \
  PORT=7000 \
  CGO_ENABLED=0 \
  GO111MODULE=on

COPY go.mod go.sum ./
RUN go mod download

COPY . .

EXPOSE $PORT

CMD ["make", "run"]

FROM dev AS build

RUN make build

FROM alpine:3.12 AS prod

WORKDIR /usr/local/bin

ARG ROOT=/go/src/github.com/noisyscanner/ivapi
COPY --from=build $ROOT/tmp/api .

CMD api
