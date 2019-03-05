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

	"github.com/dghubble/sling"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
)

type ArtifactVerifyTrackerRequest struct {
	Client   string `json:"type"`
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
