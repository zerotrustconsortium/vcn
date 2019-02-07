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
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/fatih/color"

	"../proof"
)

func verifyAll(files []string) {
	for i := 0; i < len(files); i++ {
		verify(files[i])
	}
}

func verify(filename string) {
	hash := strings.TrimSpace(hash(filename))
	verified, owner, timestamp := verifyHash(hash)
	fmt.Println("File:\t", filename)
	fmt.Println("Hash:\t", hash)
	if timestamp != 0 {
		fmt.Println("Date:\t", time.Unix(timestamp, 0))
	}
	if owner != "" {
		fmt.Println("Signer:\t", owner)
	}
	if verified {
		color.Set(color.FgHiWhite, color.BgCyan, color.Bold)
		fmt.Print("TRUSTED")
	} else {
		color.Set(color.FgHiWhite, color.BgRed, color.Bold)
		fmt.Print("UNTRUSTED")
		defer os.Exit(1)
	}
	color.Unset()
	fmt.Println()
}

func verifyHash(hash string) (verified bool, owner string, timestamp int64) {
	client, err := ethclient.Dial(MainNetEndpoint())
	if err != nil {
		log.Fatal(err)
	}
	address := common.HexToAddress(ProofContractAddress())
	instance, err := proof.NewProof(address, client)
	if err != nil {
		log.Fatal(err)
	}
	artifact, err := instance.Get(nil, hash)
	if err != nil {
		log.Fatal(err)
	}
	return artifact.Owner != "", artifact.Owner, artifact.Timestamp.Int64()
}
