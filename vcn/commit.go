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
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/dghubble/sling"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
)

type ArtifactCommitTrackerRequest struct {
	Filename string `json:"filename"`
	Hash     string `json:"hash"`
	Name     string `json:"name"`
	Url      string `json:"url"`
}

func artifactCommitTracker(hash string, filename string, status Status) {

	token, err := LoadToken()
	if err != nil {
		LOG.WithFields(logrus.Fields{
			"filename": filename,
		}).Error("Could not load token.")
	}

	// make sure the tracker does its analytics although the main
	// thread against the BC has already finalized
	defer WG.Done()

	restError := new(Error)
	r, err := sling.New().
		Post(TrackingEvent() + "/sign").
		Add("Authorization", "Bearer "+token).
		BodyJSON(ArtifactCommitTrackerRequest{
			Name:     getStatusName(int(status)),
			Hash:     hash,
			Filename: filename,
		}).Receive(nil, restError)
	if err != nil {
		fmt.Println(err)

	}
	if r.StatusCode != 200 {

		LOG.WithFields(logrus.Fields{
			"errMsg": restError.Message,
			"status": restError.Status,
		}).Error("API analytics (commit tracker) failed")
	} else {
		LOG.WithFields(logrus.Fields{
			"hash": hash,
		}).Info("Verify tracker / analytics written")
	}

}

func commitHash(hash string, passphrase string, filename string, status Status, visibility Visibility) (ret bool, code int) {
	reader, err := firstFile(WalletDirectory())
	if err != nil {
		LOG.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Could not load keystore")
	}
	transactor, err := bind.NewTransactor(reader, passphrase)
	walletSynced, err := isWalletSynced(transactor.From.Hex())
	if err != nil {
		LOG.WithFields(logrus.Fields{
			"error": err,
		}).Error("Could not load wallets")
		PrintErrorURLCustom("wallet", 400)
		os.Exit(1)
	}
	if !walletSynced {
		LOG.Error("\n", filename, " cannot be signed with CodeNotary. We are "+
			"finalizing your account configuration.\nWe will complete the "+
			"configuration shortly and we will update you as soon as this "+
			"is done.\nWe are sorry for the inconvenience and would like "+
			"to thank you for your patience.")
		os.Exit(1)
	}
	if err != nil {
		LOG.WithFields(logrus.Fields{
			"error": err,
		}).Error("Could not load contract")
		PrintErrorURLCustom("sign", 401)
		os.Exit(1)
	}
	transactor.GasLimit = GasLimit()
	transactor.GasPrice = GasPrice()
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
	tx, err := instance.Sign(transactor, hash, big.NewInt(int64(status)))
	if err != nil {
		LOG.WithFields(logrus.Fields{
			"error": err,
			"hash":  hash,
		}).Fatal("method <Sign> failed")
	}
	timeout, err := waitForTx(tx.Hash(), TxVerificationRounds(), PollInterval())
	if err != nil {
		LOG.WithFields(logrus.Fields{
			"error": err,
		}).Error("Could not write to blockchain")
		PrintErrorURLCustom("blockchain-permission", 403)
		os.Exit(1)
	}
	if timeout {
		LOG.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Writing to blockchain timed out")
	}
	publicKey, err := PublicKeyForLocalWallet()
	if err != nil {
		log.Fatal(err)
	}
	err = CreateArtifact(publicKey, filename, hash, visibility, status)
	if err != nil {
		log.Fatal(err)
	}

	return true, 0
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
