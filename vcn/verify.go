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
	"math/big"

	"github.com/dghubble/sling"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
)

type ArtifactVerifyTrackerRequest struct {
	Client   string `json:"client"`
	Filename string `json:"filename"`
	Hash     string `json:"hash"`
	Url      string `json:"url"`
}

func artifactTracker(hash string) {

	// make sure the tracker does its analytics although the main
	// thread against the BC has already finalized
	defer WG.Done()

	restError := new(Error)
	r, err := sling.New().
		Post(TrackingEvent()+"/verify").
		//Add("Authorization", "Bearer "+token).
		BodyJSON(ArtifactVerifyTrackerRequest{
			Client: "VCN:" + VCN_VERSION,
			Hash:   hash,
		}).Receive(nil, restError)
	if err != nil {
		fmt.Println(err)

	}
	if r.StatusCode != 200 {

		LOG.WithFields(logrus.Fields{
			"errMsg": restError.Message,
			"status": restError.Status,
		}).Error("API analytics failed")
	} else {
		LOG.WithFields(logrus.Fields{
			"hash": hash,
		}).Info("Verify tracker / analytics written")
	}

}

func verifyHash(hash string) (verified bool, owner string, level int64, status int64, timestamp int64) {
	client, err := ethclient.Dial(MainNetEndpoint())
	if err != nil {
		LOG.WithFields(logrus.Fields{
			"error":   err,
			"network": MainNetEndpoint(),
		}).Fatal("Cannot connect to blockchain")
	}
	address := common.HexToAddress(AssetsRelayContractAddres())
	instance, err := NewAssetsRelay(address, client)
	if err != nil {
		LOG.WithFields(logrus.Fields{
			"error":    err,
			"contract": AssetsRelayContractAddres(),
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
		l.Int64(),
		s.Int64(),
		ts.Int64()
}
