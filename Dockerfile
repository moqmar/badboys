FROM alpine:edge

RUN apk add --no-cache mariadb-client postgresql-client sqlite mongodb mongodb-tools
