# build stage
FROM golang:1.14.2-alpine3.11 AS build-env
WORKDIR /root
RUN apk --no-cache add build-base git
ADD . /root
RUN go version
RUN go build -o file-server main.go

# final stage
FROM alpine
WORKDIR /root
RUN apk --no-cache add bash
RUN mkdir -p /data/file-server
RUN mkdir -p /etc/file-server
COPY docker/file-server.json /etc/file-server/file-server.json
COPY docker/run.sh .
RUN chmod u+x run.sh
COPY --from=build-env /root/file-server /root
CMD [ "./run.sh" ]