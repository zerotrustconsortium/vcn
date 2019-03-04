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
	"encoding/json"
	"fmt"
	"log"

	"github.com/dghubble/sling"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type ArtifactVerifyTrackerRequest struct {
	Type     string `json:"type"`
	Metadata string `json:"metadata"`
}
type innerMetadata struct {
	// this is a json streucture, but seriaized to string
	Hash   string `json:"hash"`
	Client string `json:"client"`
}

func artifactTracker(hash string) {

	// make sure the tracker does its analytics although the main
	// thread against the BC has already finalized
	defer WG.Done()

	// some gymnastics with json strings
	md := &innerMetadata{Hash: hash, Client: "vcn" + VCN_VERSION}
	ser, _ := json.Marshal(md)

	restError := new(Error)
	r, err := sling.New().
		Post(TrackingEvent()).
		BodyJSON(ArtifactVerifyTrackerRequest{
			Type:     "VERIFY",
			Metadata: string(ser),
		}).Receive(nil, restError)
	if err != nil {
		fmt.Println(err)

	}
	if r.StatusCode != 200 {
		fmt.Println(fmt.Errorf("request failed: %s (%d)", restError.Message,
			restError.Status))
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
