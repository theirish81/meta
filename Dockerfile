FROM golang:1.24-alpine AS builder

RUN mkdir /usr/local/meta
WORKDIR /usr/local/meta

COPY . .

RUN go get
RUN go build -o meta *.go

FROM node:25-alpine AS web-builder
RUN mkdir /usr/local/meta
COPY . /usr/local/meta/
WORKDIR /usr/local/meta/web
RUN npm install
RUN npm run build


FROM alpine
RUN mkdir /usr/local/meta
WORKDIR /usr/local/meta
COPY --from=builder /usr/local/meta/meta .
COPY --from=web-builder /usr/local/meta/web/dist web

RUN addgroup -g 1000 meta && \
    adduser -D \
        -u 1000 \
        -G meta \
        -h /usr/local/meta \
        -s /bin/sh \
        meta

USER meta
WORKDIR /usr/local/meta
ENTRYPOINT [ "/usr/local/meta/meta" ]