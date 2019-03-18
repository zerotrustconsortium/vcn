/*
 * Copyright (c) 2018-2019 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package main

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
)

func verifyHash(hash string) (verified bool, owner string, level Level, status Status, timestamp int64) {
	client, err := ethclient.Dial(MainNetEndpoint())
	if err != nil {
		LOG.WithFields(logrus.Fields{
			"error":   err,
			"network": MainNetEndpoint(),
		}).Fatal("Cannot connect to blockchain")
	}
	address := common.HexToAddress(AssetsRelayContractAddress())
	instance, err := NewAssetsRelay(address, client)
	if err != nil {
		LOG.WithFields(logrus.Fields{
			"error":    err,
			"contract": AssetsRelayContractAddress(),
		}).Fatal("Cannot instantiate contract")
	}
	address, l, s, ts, err := instance.Verify(nil, hash)
	if err != nil {
		LOG.WithFields(logrus.Fields{
			"error": err,
			"hash":  hash,
		}).Fatal("method <Verify> failed")
	}
	return address != common.BigToAddress(big.NewInt(0)),
		address.Hex(),
		Level(l.Int64()),
		Status(s.Int64()),
		ts.Int64()
}
