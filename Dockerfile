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

ENV XDEBUG 9003:Development

ENV IDE_IP host.docker.internal
ENV IDE_PORT 9010

ENV FRAMEWORK "flow"

ENV LOCAL_ROOT ""

ENTRYPOINT ["sh", "-c", "./flow-debugproxy --xdebug ${XDEBUG} --framework ${FRAMEWORK} --ide ${IDE_IP}:${IDE_PORT} --localroot \"${LOCAL_ROOT}\" ${ADDITIONAL_ARGS}"]
