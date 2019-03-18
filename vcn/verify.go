/*
 * Copyright (c) 2018-2019 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package main

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type BlockchainVerification struct {
	Owner     common.Address
	Level     Level
	Status    Status
	Timestamp time.Time
}

func VerifyHash(hash string) (verification *BlockchainVerification, err error) {
	client, err := ethclient.Dial(MainNetEndpoint())
	if err != nil {
		return nil, err
	}
	contractAddress := common.HexToAddress(AssetsRelayContractAddress())
	instance, err := NewAssetsRelay(contractAddress, client)
	if err != nil {
		return nil, err
	}
	address, level, status, timestamp, err := instance.Verify(nil, hash)
	if err != nil {
		return nil, err
	}
	verification = new(BlockchainVerification)
	verification.Owner = address
	verification.Level = Level(level.Int64())
	verification.Status = Status(status.Int64())
	verification.Timestamp = time.Unix(timestamp.Int64(), 0)
	return verification, nil
}
