FROM alpine

WORKDIR /go/tstl

ADD . ./

RUN chmod +x /go/tstl/start.sh

EXPOSE 10080

CMD ["./tstl"]