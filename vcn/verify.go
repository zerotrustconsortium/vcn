/*
 * Copyright (c) 2018-2019 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 * Built on top of CLI (MIT license)
 * https://github.com/urfave/cli#overview
 */

package main

import (
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func verifyHash(hash string) (verified bool, owner string, timestamp int64) {
	client, err := ethclient.Dial(MainNetEndpoint())
	if err != nil {
		log.Fatal(err)
	}
	address := common.HexToAddress(ProofContractAddress())
	instance, err := NewProof(address, client)
	if err != nil {
		log.Fatal(err)
	}
	artifact, err := instance.Get(nil, hash)
	if err != nil {
		log.Fatal(err)
	}
	return artifact.Owner != "", artifact.Owner, artifact.Timestamp.Int64()
}
