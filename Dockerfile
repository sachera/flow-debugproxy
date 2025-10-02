#
# Compile step
FROM golang:alpine AS build-env
ENV GOPATH=/gopath
ENV PATH=$GOPATH/bin:$PATH
WORKDIR /app
RUN apk update && \
    apk upgrade && \
    apk add git
ADD go.mod .
ADD go.sum .
RUN go mod download
ADD . .
RUN CGO_ENABLED=0 go build -o flow-debugproxy

#
# Build step
FROM alpine
WORKDIR /app

COPY --from=build-env /app/flow-debugproxy /app/

ENV ADDITIONAL_ARGS ""

ENV XDEBUG_PORT 9010

ENV IDE_IP 127.0.0.1
ENV IDE_PORT 9000

ENV FRAMEWORK "flow"

ENTRYPOINT ["sh", "-c", "./flow-debugproxy --xdebug 0.0.0.0:${XDEBUG_PORT} --framework ${FRAMEWORK} --ide ${IDE_IP}:${IDE_PORT} ${ADDITIONAL_ARGS}"]
