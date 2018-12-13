#!/bin/sh
set -eux

# if first arg is node in monorepo graph - apply commands to this node (package),
# otherwise - to root (all repo nodes)
NODE=${1:-.}

go get -d -t ./${NODE}/...
go test -race -run=^$$ -bench=. -benchmem ./${NODE}/...
