# multi stage build
FROM golang:1.21.4 as builder
WORKDIR /go/setenv
ENV CGO_ENABLED=0
RUN mkdir setenv && \
    apt-get update && \
    apt-get -y install vim && \
    apt-get -y install unzip
RUN curl -OL https://github.com/google/protobuf/releases/download/v3.11.2/protoc-3.11.2-linux-x86_64.zip && \
    unzip protoc-3.11.2-linux-x86_64.zip -d protoc3 && \
    mv protoc3/bin/* /usr/local/bin/ && \
    mv protoc3/include/* /usr/local/include/ && \
    rm -rf protoc-3.11.2-linux-x86_64.zip protoc3
COPY . .
RUN go build .

# As runner
FROM alpine:3.18.4
RUN apk --no-cache add ca-certificates &&\
    apk add gcompat
COPY --from=builder  /go/setenv/recommend .
CMD ["./recommend"]