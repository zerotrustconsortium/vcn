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
	"context"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func commitHash(hash string, owner string, passphrase string, filename string) {
	reader, err := firstFile(WalletDirectory())
	if err != nil {
		log.Fatal(err)
	}
	transactor, err := bind.NewTransactor(reader, passphrase)
	if err != nil {
		log.Fatal(err)
	}
	transactor.GasLimit = GasLimit()
	transactor.GasPrice = GasPrice()
	client, err := ethclient.Dial(MainNetEndpoint())
	if err != nil {
		log.Fatal(err)
	}
	address := common.HexToAddress(ProofContractAddress())
	instance, err := NewProof(address, client)
	if err != nil {
		log.Fatal(err)
	}
	tx, err := instance.Set(transactor, owner, hash)
	if err != nil {
		log.Fatal(err)
	}
	timeout, err := waitForTx(tx.Hash(), TxVerificationRounds(), PollInterval())
	if err != nil {
		log.Fatal(err)
	}
	if timeout {
		log.Fatal("transaction timed out")
	}
	publicKey, err := PublicKeyForLocalWallet()
	if err != nil {
		log.Fatal(err)
	}
	err = CreateArtifact(publicKey, filename, hash)
	if err != nil {
		log.Fatal(err)
	}
}

func waitForTx(tx common.Hash, maxRounds uint64, pollInterval time.Duration) (timeout bool, err error) {
	client, err := ethclient.Dial(MainNetEndpoint())
	if err != nil {
		return false, err
	}
	for i := uint64(0); i < maxRounds; i++ {
		_, pending, err := client.TransactionByHash(context.Background(), tx)
		if err != nil {
			return false, err
		}
		if !pending {
			return false, nil
		}
		time.Sleep(pollInterval)
	}
	return true, nil
}
