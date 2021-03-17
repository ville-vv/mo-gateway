FROM alpine

WORKDIR /go/mo-gateway

ADD . ./

RUN chmod +x /go/github.com/ville-vv/mo-gateway/start.sh

EXPOSE 10080

CMD ["./mo-gateway"]