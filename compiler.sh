#!/bin/bash
compiler() {
  CGO_ENABLED=0 GOOS=$1 GOARCH=$2 go build -trimpath -o bin/dsvha-${1}-${2} main.go
}
compiler linux amd64
compiler linux 386
compiler linux arm64
compiler linux arm
compiler windows amd64
compiler windows 386
compiler darwin amd64
#compiler darwin 386