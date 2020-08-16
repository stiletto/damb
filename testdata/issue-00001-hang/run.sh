#!/bin/bash
set -eu -o pipefail
test -f nope/random-large-file.bin || dd if=/dev/zero of=nope/random-large-file.bin bs=1M count=1
damb build nope && exit 1 || exit 0
