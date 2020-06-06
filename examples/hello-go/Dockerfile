ARG shit=docker.io
FROM ${shit}/golang:1 AS builder
WORKDIR /root
ADD hello.go /root/
RUN go build -o hello

FROM builder as fuck

FROM ${shit}/busybox:1-glibc
COPY --from=builder /root/hello /bin/hello

ENTRYPOINT ["/bin/hello"]
