FROM golang:1.16-alpine AS mybuildstage

RUN apk add build-base
RUN apk add sqlite

ADD project-src /usr/project

RUN sqlite3 sqlite.db "VACUUM;"

WORKDIR /usr/project

RUN go build -o app src/main.go

FROM alpine:3.12.2

COPY --from=mybuildstage /usr/project /usr/project

RUN /bin/sh -c "apk add --no-cache bash"