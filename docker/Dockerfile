FROM alpine

WORKDIR /go/mo-gateway

ADD . ./

RUN chmod +x /go/mo-gateway/start.sh

EXPOSE 10080

CMD ["./mo-gateway"]