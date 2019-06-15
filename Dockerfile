FROM golang:1.12.5
WORKDIR /go/src/github.com/leosunmo/simple-redirector/
#RUN go get -d -v golang.org/x/net/html  
COPY main.go .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo .

FROM alpine:3.9  
#RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /go/src/github.com/leosunmo/simple-redirector/ .
CMD ["./simple-redirector"]  