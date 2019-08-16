FROM golang:1.12.8-alpine
WORKDIR /go/src/github.com/mayueeeee/KMITL-Auto-Authen
COPY . .
RUN apk add --no-cache git \
    && go get github.com/joho/godotenv \
    && apk del git \
    && go build auth.go

FROM alpine:3.10.1
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /go/src/github.com/mayueeeee/KMITL-Auto-Authen .
CMD ["./app"]  
