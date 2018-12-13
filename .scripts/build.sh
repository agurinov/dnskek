#!/bin/sh

<< ////

TODO:
use cross-compiling
$go get github.com/mitchellh/gox
$gox \
	-verbose \
	-osarch="linux/amd64 darwin/amd64" \
	-ldflags="-X main.VERSION=$CIRCLE_TAG -X main.TIMESTAMP=`date +%s`" \
	-output="/go/bin/{{.Dir}}-{{.OS}}-{{.Arch}}"

////

set -eux

# if first arg is node in monorepo graph - apply commands to this node (package),
# otherwise - to root (all repo nodes)
NODE=${1:-.}

go get -d ./${NODE}/...

# calculate base variables
BASE=`basename "${NODE}"`
TIMESTAMP=`date +%s`
VERSION="${CIRCLE_TAG:=LOCAL}"

# calculate build/compile specific variables
ldflags="-X 'main.VERSION=${VERSION}' -X 'main.TIMESTAMP=${TIMESTAMP}'"

# compile and build
# linux 64 bit
# - alpine, ubuntu
GOOS=linux GOARCH=amd64 go build \
	-v \
	-ldflags "${ldflags}" \
	-o /go/bin/${BASE}-Linux-x86_64 \
	./${NODE} # otherwise -> go build: cannot use -o with multiple packages

# - macos
GOOS=darwin GOARCH=amd64 go build \
	-v \
	-ldflags "${ldflags}" \
	-o /go/bin/${BASE}-Darwin-x86_64 \
	./${NODE} # otherwise -> go build: cannot use -o with multiple packages

# chmod binaries to executable
chmod +x /go/bin/${BASE}-Linux-x86_64
chmod +x /go/bin/${BASE}-Darwin-x86_64
