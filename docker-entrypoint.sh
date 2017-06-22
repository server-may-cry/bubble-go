#!/bin/bash
set -e

make deps
make build

exec "$@"
