/*
 * Copyright (c) 2018-2019 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 *
 * User Interaction
 *
 * This part of the vcn code handles the concern of interaction (the *V*iew)
 *
 */

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"

	"github.com/pkg/browser"
	"golang.org/x/crypto/ssh/terminal"
)

func dashboard() {
	// open dashboard
	// we intentionally do not read the customer's token from disk
	// and GET the dashboard => this would be insecure as tokens would
	// be visible in server logs. in case the anyhow long-running web session
	// has expired the customer will have to log in
	url := DashboardURL()
	fmt.Println(fmt.Sprintf("Taking you to <%s>", url))
	browser.OpenURL(url)
}

func auth() {

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Email address: ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSuffix(email, "\n")

	fmt.Print("     Password: ")
	password, err := terminal.ReadPassword(int(syscall.Stdin))
	passwordString := string(password)
	fmt.Println("")

	if err != nil {
		log.Fatal(err)
	}

	ret := Authenticate(email, passwordString)

	if ret {
		fmt.Println("Authentication successful.")
	}

	// check for a keystore right now and hint at creating one
	// this is for the case of a newly registered customer
	// coming in from the dashboard
	// doing vcn auth for the very first time
	// and no keystore is yet present
	hasKeystore, err := HasKeystore()
	if err != nil {
		log.Fatal(err)
	}
	if hasKeystore == false {
		fmt.Println("You have no keystore set up yet.")
		fmt.Println("<vcn> will now do this for you and upload the public key to the platform.")

		keystorePassphrase, err := readPassword("Keystore passphrase:")
		if err != nil {
			log.Fatal(err)
		}

		CreateKeystore(keystorePassphrase)
		SyncKeys()

	}

}
