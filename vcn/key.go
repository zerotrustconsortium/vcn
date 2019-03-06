/*
 * Copyright (c) 2018-2019 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/dghubble/sling"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/sirupsen/logrus"
)

type Wallet struct {
	Id      uint64 `json:"id"`
	Address string `json:"address"`
}

type PagedWalletResponse struct {
	Content []Wallet `json:"content"`
}

func CreateKeystore(password string) (pubKey string, wallet string) {
	if password == "" {
		LOG.Error("Keystore passphrase cannot be empty")
	}
	ks := keystore.NewKeyStore(WalletDirectory(), keystore.StandardScryptN, keystore.StandardScryptP)
	account, err := ks.NewAccount(password)
	if err != nil {
		log.Fatal(err)
	}

	pubKey = account.Address.Hex()
	wallet = WalletDirectory()

	return pubKey, wallet
}

func HasKeystore() (bool, error) {

	LOG.WithFields(logrus.Fields{
		"keystore": WalletDirectory(),
	}).Trace("HasKeystore()")

	files, err := ioutil.ReadDir(WalletDirectory())
	if err != nil {
		LOG.WithFields(logrus.Fields{
			"error": err,
		}).Error("ReadDir() failed")
		return false, err
	}
	return len(files) > 0, nil
}

func LoadPublicKeys() (addresses []string, err error) {
	authError := new(Error)
	pagedWalletResponse := new(PagedWalletResponse)
	token, err := LoadToken()
	if err != nil {
		log.Fatal(err)
	}
	r, err := sling.New().
		Add("Authorization", "Bearer "+token).
		Get(WalletEndpoint()).
		Receive(pagedWalletResponse, authError)
	if err != nil {
		log.Fatal(err)
	}
	if r.StatusCode != 200 {
		log.Fatalf("request failed: %s (%d)", authError.Message,
			authError.Status)
	}
	var result []string
	for _, wallet := range (*pagedWalletResponse).Content {
		result = append(result, wallet.Address)
	}
	return result, nil
}

func SyncKeys() {
	authError := new(Error)
	token, err := LoadToken()
	if err != nil {
		log.Fatal(err)
	}
	addresses, err := LoadPublicKeys()
	if err != nil {
		log.Fatal(err)
	}
	localAddress, err := PublicKeyForLocalWallet()
	if err != nil {
		log.Fatal(err)
	}
	if contains(addresses, localAddress) {
		return
	}
	r, err := sling.New().
		Add("Authorization", "Bearer "+token).
		Post(WalletEndpoint()).
		BodyJSON(Wallet{Address: localAddress}).
		Receive(nil, authError)
	if err != nil {
		log.Fatal(err)
	}
	if r.StatusCode != 200 {
		log.Fatalf("request failed: %s (%d)", authError.Message,
			authError.Status)
	}

	// TODO: tracking key added
	//fmt.Println(">>>>>>>>>> WILL ADD KEY ADDED HERE")

}

func PublicKeyForLocalWallet() (string, error) {
	reader, err := firstFile(WalletDirectory())
	if err != nil {
		return "", err
	}
	contents, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}
	var keyfile map[string]*json.RawMessage
	err = json.Unmarshal(contents, &keyfile)
	if err != nil {
		return "", err
	}
	var localAddress string
	err = json.Unmarshal(*keyfile["address"], &localAddress)
	if err != nil {
		return "", err
	}
	return "0x" + localAddress, nil
}
