FROM golang:1.23.0-alpine3.19 AS build
RUN apk --no-cache add gcc g++ make ca-certificates
WORKDIR /github.com/donaldnash/go-marketplace
COPY go.mod go.sum ./
COPY account account
RUN go build -o /go/bin/app ./account/cmd/account

FROM alpine:3.19
WORKDIR /usr/bin
COPY --from=build /go/bin .
EXPOSE 8080
CMD ["app"]