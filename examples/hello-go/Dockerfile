ARG foo=docker.io
FROM ${foo}/golang:1 AS builder
WORKDIR /root
ADD hello.go /root/
RUN go build -o hello

FROM builder as fuck

FROM ${foo}/busybox:1-glibc
COPY --from=builder /root/hello /bin/hello

ENTRYPOINT ["/bin/hello"]
