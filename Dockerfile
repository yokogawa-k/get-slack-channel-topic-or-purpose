FROM golang:1.12 as Build

ENV GO111MODULE=on
WORKDIR /go/src/app
COPY go.mod go.sum ./

RUN set -ex \
    && go mod download

COPY ./ ./

RUN set -ex \
    && go build -o get-slacklog .

FROM scratch

COPY --from=Build /go/src/app/get-slacklog /
