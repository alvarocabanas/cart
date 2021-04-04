FROM golang:1.14-alpine as builder

ENV CGO_ENABLED=0
ENV GOBIN=/go/bin

WORKDIR /go/src/cart
COPY ../../../.. /go/src/cart
RUN cd /go/src/cart && go install -a -ldflags '-w' cmd/server/main.go

FROM alpine:latest
COPY --from=builder /go/bin/main /usr/bin/cmd
CMD ["/usr/bin/cmd"]