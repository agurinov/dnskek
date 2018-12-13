#!/bin/sh
set -eux

# if first arg is node in monorepo graph - apply commands to this node (package),
# otherwise - to root (all repo nodes)
NODE=${1:-.}

go get -d -t ./${NODE}/...
go clean -testcache || true
go test -race -cover ./${NODE}/...

# other variants:
# go test  -v -race -cover ./...
# go test -v -race -cover ./...
# go test -race -cover tools/context_test.go
