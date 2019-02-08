/*
 * Copyright (c) 2018-2019 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package main

import (
	"fmt"
	"log"
	"syscall"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"golang.org/x/crypto/ssh/terminal"
)

func createKs() {
	ks := keystore.NewKeyStore(WalletDirectory(), keystore.StandardScryptN, keystore.StandardScryptP)
	fmt.Print("Keystore password:")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	password := string(bytePassword)
	fmt.Println(".")
	account, err := ks.NewAccount(password)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Public key:\t", account.Address.Hex())
	fmt.Println("Keystore:\t", WalletDirectory())
}
