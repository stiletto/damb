#!/bin/sh
echo "Hello from busybox"
exec /bin/hello "$@"
