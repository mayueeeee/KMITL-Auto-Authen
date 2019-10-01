FROM golang:latest

WORKDIR /auth

COPY ./ /auth/

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o auth .

ENTRYPOINT ["./auth"]