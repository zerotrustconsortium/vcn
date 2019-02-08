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
	"github.com/fatih/color"
	"golang.org/x/crypto/ssh/terminal"
)

// https://github.com/miguelmota/ethereum-development-with-go-book/blob/master/en/keystore/README.md
func createKs() {

	ks := keystore.NewKeyStore(WalletDirectory(), keystore.StandardScryptN, keystore.StandardScryptP)

	fmt.Print("Choose a password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	password := string(bytePassword)
	fmt.Println(".")

	account, err := ks.NewAccount(password)
	if err != nil {
		log.Fatal(err)
	}

	output(account.Address.Hex())

	fmt.Println("I've also put it to", WalletDirectory())

}
func output(addr string) {
	fmt.Print("Okay, that's your public key: ")
	color.Set(color.FgHiWhite, color.BgCyan, color.Bold)

	fmt.Printf("%s", addr)
	color.Unset()
	fmt.Println("")
}
