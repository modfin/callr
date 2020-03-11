#!/usr/bin/env bash

VERSION=$(date +%Y-%m-%dT%H.%M.%S)-$(git log -1 --pretty=format:"%h")
IMAGE_NAME=modfin/callrd

docker build -f ./cmd/callrd/Dockerfile \
    -t ${IMAGE_NAME}:latest \
    -t ${IMAGE_NAME}:${VERSION} \
    . || exit 1
