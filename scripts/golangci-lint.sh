#!/bin/env bash

VERSION=v1.33.0
PROJECT_ROOT=$(dirname $(dirname $(realpath ${BASH_SOURCE[0]})))

exec docker run --rm -v $PROJECT_ROOT:/app -w /app golangci/golangci-lint:$VERSION \
    golangci-lint run -v
