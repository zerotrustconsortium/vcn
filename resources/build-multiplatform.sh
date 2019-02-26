#!/bin/bash
# NAIVE approach does not work with C libraries of go-ethereum
#env GOOS=windows GOARCH=amd64 go build
#env GOOS=linux GOARCH=amd64 go build
read -p 'Version (9.9.9): ' version
mkdir -p ../dist
xgo --targets=linux/amd64,windows/amd64,darwin/amd64 -out ../dist/vcn-v$version ../vcn
