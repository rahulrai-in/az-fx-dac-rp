
FROM golang:1.15.2-alpine AS builder

RUN apk update && apk add git && apk add ca-certificates

WORKDIR /az-fx-proxy

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/az-fx-proxy

# Runtime image
FROM scratch AS base
COPY --from=builder /etc/ssl/certs/ /etc/ssl/certs/
COPY --from=builder /go/bin/az-fx-proxy /bin/az-fx-proxy
ENTRYPOINT ["/bin/az-fx-proxy"]