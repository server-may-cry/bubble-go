#!/bin/sh
set -e

make deps

exec "$@"
