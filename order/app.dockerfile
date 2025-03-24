FROM golang:1.22-alpine3.19 AS build
RUN apk --no-cache add gcc g++ make ca-certificates
WORKDIR /github.com/donaldnash/go-marketplace
COPY go.mod go.sum ./
COPY vendor vendor
COPY account account
COPY catalog catalog
COPY order order
RUN go build -o /go/bin/app ./order/cmd/order

FROM alpine:3.19
WORKDIR /usr/bin
COPY --from=build /go/bin .
EXPOSE 8080
CMD ["app"]