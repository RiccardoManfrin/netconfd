# syntax=docker/dockerfile:experimental
FROM golang:1.16.2-alpine3.13 as builder
ARG CI_PROJECT_PATH
RUN mkdir -p -m 0600 ~/.ssh
RUN echo "StrictHostKeyChecking no" >> ~/.ssh/config
RUN apk update && \
    apk add --no-cache git make openssh pkgconfig gcc
RUN update-ca-certificates
ENV GOPATH /go/
ENV SRCPATH $GOPATH/src/github.com/$CI_PROJECT_PATH
ADD src.tar.gz $SRCPATH
RUN --mount=type=ssh \
	cd $SRCPATH && \
    make
##coverage
#RUN apk add libc-dev
#RUN cd $SRCPATH && \
#    go test -cover

FROM alpine:3.7
ARG CI_PROJECT_PATH
ARG CI_PROJECT_NAME
ENV CI_PROJECT_NAME=$CI_PROJECT_NAME
ENV GOPATH /go/
ENV SRCPATH $GOPATH/src/gitlab.lan.athonet.com/$CI_PROJECT_PATH
RUN mkdir /app/
COPY --from=builder $SRCPATH/$CI_PROJECT_NAME /app/$CI_PROJECT_NAME
RUN mkdir /conf/
COPY --from=builder $SRCPATH/$CI_PROJECT_NAME.json /conf/$CI_PROJECT_NAME.json.golden
WORKDIR /app/
ENTRYPOINT ./$CI_PROJECT_NAME --config /conf/$CI_PROJECT_NAME.json
