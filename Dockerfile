FROM golang:alpine as builder
RUN mkdir -p $GOPATH/src/github.com/valerykalashnikov/moments
ADD . $GOPATH/src/github.com/valerykalashnikov/moments
RUN cd $GOPATH/src && go build github.com/valerykalashnikov/moments/example/http_server 

FROM alpine
WORKDIR /app
COPY --from=builder /go/src/http_server /app/
ENTRYPOINT ./http_server