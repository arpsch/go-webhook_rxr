FROM golang:1.16-alpine as build

# Copy the code from the host and compile it
WORKDIR go-webhook_rxr
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o wh_receiver

FROM alpine:3.6

EXPOSE 9999

COPY --from=build ./go/go-webhook_rxr/wh_receiver /usr/bin/

ENTRYPOINT ["/usr/bin/wh_receiver"]


