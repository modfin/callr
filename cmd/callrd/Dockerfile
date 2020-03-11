FROM golang:1.14-alpine as builder
RUN apk add --no-cache tzdata ca-certificates

RUN mkdir -p /go/src/callr
COPY . /go/src/callr
RUN cd /go/src/callr && go build -o /callrd ./cmd/callrd/*.go

FROM alpine
RUN apk add --no-cache tzdata ca-certificates
EXPOSE 8080
COPY --from=builder /callrd /
CMD /callrd