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
	"time"

	"github.com/fatih/color"
	"github.com/pkg/browser"
	"golang.org/x/crypto/ssh/terminal"
)

func register() {

	var keystorePassphrase string

	hasKeystore, err := HasKeystore()
	if err != nil {
		log.Fatal(err)
	}
	if hasKeystore {
		// for now, we only enable ONE keystore
		fmt.Printf("You already have a keystore.\n")
		PrintErrorURLCustom("keystore", 428)
		os.Exit(1)
	}

	fmt.Println("User registration for vChain.us\n" +
		"If you already have an account abort and\n" +
		"authenticate yourself using <vcn auth>")

	email, accountPassword, _ := getLoginCredentials()

	Register(email, accountPassword)

	if !hasKeystore {
		keystorePassphrase, err = readPassword("Keystore passphrase:")
		if err != nil {
			log.Fatal(err)
		}

		pubKey, wallet := CreateKeystore(keystorePassphrase)

		fmt.Println("Keystore successfully created")
		fmt.Println("Public key:\t", pubKey)
		fmt.Println("Keystore:\t", wallet)

		fmt.Println("We've sent you an email to: ", email,
			"\nClick the link and you will be automatically logged in.")
		color.Set(StyleAffordance())
		fmt.Print("Check your email [...]")
		color.Unset()
		fmt.Println()

		err = WaitForConfirmation(email, accountPassword,
			60, 2*time.Second)
		if err != nil {
			log.Fatal(err)
		}
	}

	SyncKeys()
}

func getLoginCredentials() (email string, passwordString string, err error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Email address: ")
	email, _ = reader.ReadString('\n')
	email = strings.TrimSuffix(email, "\n")

	fmt.Print("Password: ")
	password, err := terminal.ReadPassword(int(syscall.Stdin))
	passwordString = string(password)

	fmt.Println("")

	if err != nil {
		log.Fatal(err)
	}

	return email, passwordString, nil
}

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

	} else {
		SyncKeys()
	}

}

// Commit => "sign"
func Sign(filename string, owner string) {
	hash := hash(filename)

	reader := bufio.NewReader(os.Stdin)

	fmt.Println()

	output := "" +
		"vChain CodeNotary - code signing made easy:\n" +
		"-------------------------------------------\n" +
		"Attention, by signing this artifact you implicitly claim its ownership.\n" +
		"Doing this can potentially infringe other publisher's intellectual\n" +
		"property under the laws of your country of residence.\n" +
		"vChain, CodeNotary and the Zero Trust Consortium cannot be\n" +
		"held responsible for legal ramifications.\n\n"

	fmt.Println(output)

	color.Set(color.FgGreen)
	fmt.Println("It's safe to continue if you are the owner of the artif,\ne.g. author, creator, publisher.")
	color.Unset()
	fmt.Print("\nDo you understand and want to continue? (y/N):")

	question, _ := reader.ReadString('\n')

	if strings.TrimSpace(question) != "y" {

		fmt.Println("Ok - exiting.")
		os.Exit(0)
	}

	fmt.Print("Keystore passphrase:")
	passphrase, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Println(".")
	if err != nil {
		log.Fatal(err)
	}

	go displayLatency()

	// TODO: return and display: block #, trx #
	commitHash(hash, owner, string(passphrase), filename)
	fmt.Println("")
	fmt.Println("File:\t", filename)
	fmt.Println("Hash:\t", hash)
	fmt.Println("Date:\t", time.Now())
	fmt.Println("Signer:\t", owner)
}

// VerifyAll => main entry point from cli
// unwraps potential list input for verify()
func VerifyAll(files []string) {
	for _, file := range files {
		verify(file)
	}
}
func verify(filename string) {
	hash := strings.TrimSpace(hash(filename))
	verified, owner, timestamp := verifyHash(hash)
	fmt.Println("File:\t", filename)
	fmt.Println("Hash:\t", hash)
	if timestamp != 0 {
		fmt.Println("Date:\t", time.Unix(timestamp, 0))
	}
	if owner != "" {
		fmt.Println("Signer:\t", owner)
	}

	fmt.Print("Trust:\t")
	if verified {
		color.Set(color.FgHiWhite, color.BgCyan, color.Bold)
		fmt.Print("VERIFIED")
	} else {
		color.Set(color.FgHiWhite, color.BgMagenta, color.Bold)
		fmt.Print("UNKNOWN")
		defer os.Exit(1)
	}
	color.Unset()
	fmt.Println()
}

func displayLatency() {
	i := 0
	for {
		i++
		fmt.Printf("\033[2K\rIn progress %02dsec", i)
		// fmt.Println(i)
		time.Sleep(1 * time.Second)
	}
}
