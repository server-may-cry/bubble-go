#!/bin/sh
set -e

make deps
make build

exec "$@"
