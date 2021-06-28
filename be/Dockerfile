FROM alpine
RUN apk add --no-cache -U curl ca-certificates fuse && rm -rf /tmp/*
RUN mkdir -p /usr/local/go/lib/time/
RUN curl -L -o /usr/local/go/lib/time/zoneinfo.zip https://github.com/golang/go/raw/master/lib/time/zoneinfo.zip
COPY ./bin/* /usr/bin/
WORKDIR /ws/olvrng.rbot
#USER root

CMD ["sh"]
EXPOSE 8080 8180 9001 9002 9011 9012
