FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git
RUN apk --no-cache add ca-certificates

WORKDIR /src
COPY ./go.mod ./go.sum ./
RUN go mod download

COPY ./ ./

RUN go get -d -v
RUN CGO_ENABLED=0 go build -o /go/bin/dyndns-cloudflare

FROM scratch
COPY --from=builder /etc/ssl/certs/ /etc/ssl/certs/
COPY --from=builder /go/bin/dyndns-cloudflare /go/bin/dyndns-cloudflare

EXPOSE 8080

ENTRYPOINT ["/go/bin/dyndns-cloudflare"]