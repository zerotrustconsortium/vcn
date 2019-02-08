#!/bin/bash -e

docker run -v "$(pwd):/tmp" ethereum/solc:0.4.24 --abi /tmp/resources/Proof.sol -o /tmp
docker run -v "$(pwd):/tmp" ethereum/client-go:alltools-stable abigen --abi /tmp/Proof.abi --pkg main --type Proof --out /tmp/vcn/proof.go
rm -f Proof.abi
