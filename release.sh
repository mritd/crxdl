#!/bin/bash

set -e

if [ ! -n "$1" ]; then
    echo "Error:release version is blank!"
	exit 1
fi

if [ -d dist ]; then
    rm -rf dist/*
else
    mkdir dist
fi

GOOS=darwin packr build && mv crxdl dist/crxdl_darwin_amd64
GOOS=linux packr build && mv crxdl dist/crxdl_linux_amd64

ghr -u mritd -t ${GITHUB_RELEASE_TOKEN} -replace -recreate --debug $1 dist/
