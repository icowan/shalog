FROM golang:1.13.4-alpine3.10 as build-env

ENV GO111MODULE=on
ENV GOPROXY=http://goproxy.cn
ENV BUILDPATH=github.com/icowan/shalom
RUN mkdir -p /go/src/${BUILDPATH}
COPY ./ /go/src/${BUILDPATH}
RUN cd /go/src/${BUILDPATH} && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install -v

FROM alpine:latest

COPY --from=build-env /go/bin/shalom /go/bin/shalom
COPY ./views /go/bin/views
COPY ./database /go/bin/database

WORKDIR /go/bin/
CMD ["/go/bin/shalom", "start", "-p", ":8080", "-c", "/etc/shalom/app.cfg"]
