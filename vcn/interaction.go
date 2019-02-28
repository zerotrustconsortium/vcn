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

func login() {
	// login handles both first-time registration of users
	// and subsequent logins, unfortunately a
	// rather complicated combination of cases
	// ---------------------------------------
	// filesystem: token exists && api: valid
	// no  => enter email
	//        api: publisher exists
	//        yes => enter pw
	//               authenticate()
	//               fails => retry pw entry up to 3 times
	//        no  => display: registration texts
	//               confirmation that email is correct (otherwise typo leads to registration)
	//               filesystem: keystore exists
	//               yes => error "cannot register new customer with existing keystore"
	//               enter password
	//               api: register user
	//               api: authenticate()
	// api: check if verified customer
	// no  => point to verification mail
	// filesystem: keystore exists
	// no => createkeystore
	// synckeys

	token, _ := LoadToken()
	tokenValid := CheckToken(token)

	if tokenValid == false {
		// most probably a new customer or an expired token

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Email address: ")
		email, _ := reader.ReadString('\n')
		email = strings.TrimSuffix(email, "\n")

		publisherExists := CheckPublisherExists(email)

		if publisherExists {

			authenticated := false
			counter := 0

			for authenticated == false {

				counter++
				attempt := ""
				if counter == 2 {
					attempt = " (next try)"
				} else if counter == 3 {
					attempt = " (final try)"
				} else if counter == 4 {
					PrintErrorURLCustom("password", 404)
					os.Exit(1)
				}

				fmt.Printf("Password%s: ", attempt)
				password, _ := terminal.ReadPassword(int(syscall.Stdin))
				passwordString := string(password)
				fmt.Println("")

				returnCode := 0
				authenticated, returnCode = Authenticate(email, passwordString)
				if returnCode > 0 {

					if returnCode == 401 {
						fmt.Println("Please enter a correct password.")

					}

				}

			}
		} else {
			// obviously this is a new registration
			fmt.Println("vChain.us Publisher Registration")
			color.Set(StyleAffordance())
			fmt.Print("Attention: vcn will create a new account for you.")
			color.Unset()
			fmt.Printf("\nIs this your email address <%s>? (y/N):", email)

			question, _ := reader.ReadString('\n')

			if strings.TrimSpace(question) != "y" {
				fmt.Println("Ok - exiting.")
				os.Exit(1)
			}

			hasKeystore, _ := HasKeystore()
			if hasKeystore == true {

				fmt.Println("Cannot register customer with existing keystore.")
				PrintErrorURLCustom("keystore", 428)
				os.Exit(1)
			}

			fmt.Print("Choose an account password: ")
			password, _ := terminal.ReadPassword(int(syscall.Stdin))
			accountPassword := string(password)
			fmt.Println("")

			ret, returnCode := Register(email, accountPassword)
			if ret == false {
				// somehow hacky to put this into interaction.go
				// nonetheless Register is now a clean interface void of cli I/O
				PrintErrorURLByEndpoint(PublisherEndpoint(), "POST", returnCode)
				os.Exit(1)
			}

			authenticated, _ := Authenticate(email, accountPassword)
			if authenticated == false {
				fmt.Println("Could not log you in.")
				os.Exit(1)
			}

		}

	}
	// api: check if verified customer
	token, _ = LoadToken() // looks superfluous; would not be set if auth / register only happens above
	verified, _ := CheckPublisherIsVerified(token)

	if verified == false {
		fmt.Println("We've sent you an email to the address you provided." +
			"\nClick the link and you will be automatically logged in.")
		color.Set(StyleAffordance())
		fmt.Print("Check your email [...]")
		color.Unset()
		fmt.Println()
		err := WaitForConfirmation(token,
			60, 2*time.Second)
		if err != nil {
			log.Fatal(err)
		}
	}
	// filesystem: keystore exists
	// no => createkeystore
	// synckeys
	hasKeystore, err := HasKeystore()
	if err != nil {
		log.Fatal(err)
	}
	if hasKeystore == false {

		fmt.Println("You have no keystore set up yet.")
		fmt.Println("<vcn> will now do this for you and upload the public key to the platform.")

		color.Set(StyleAffordance())
		fmt.Print("Attention: Please pick a strong password. There is no recovery possible.")
		color.Unset()
		fmt.Println()

		keystorePassphrase, _ := readPassword("Keystore passphrase: ")

		pubKey, wallet := CreateKeystore(keystorePassphrase)

		fmt.Println("Keystore successfully created")
		fmt.Println("Public key:\t", pubKey)
		fmt.Println("Keystore:\t", wallet)

	}
	SyncKeys()

	fmt.Println("Login successful.")
}

// Commit => "sign"
func Sign(filename string, owner string) {

	// check for token
	token, _ := LoadToken()
	if token == "" {
		fmt.Println("You need to be logged in to sign.")
		fmt.Println("Proceed by authenticating yourself using <vcn auth>")
		//PrintErrorURLCustom("token", 428)
		os.Exit(1)
	}

	// keystore
	hasKeystore, _ := HasKeystore()
	if hasKeystore == false {
		fmt.Printf("You need a keystore to sign.\n")
		fmt.Println("Proceed by authenticating yourself using <vcn auth>")
		//PrintErrorURLCustom("keystore", 428)
		os.Exit(1)
	}

	var artifactHash string

	// artifact types
	if strings.HasPrefix(filename, "docker:") {

		artifactHash = getDockerHash(filename)
		//fmt.Printf("Docker: Not yet implemented\n")
		//os.Exit(1)

	} else if strings.HasPrefix(filename, "git:") {
		fmt.Printf("git: Not yet implemented\n")
		os.Exit(1)
	} else {
		// file mode
		artifactHash = hash(filename)

	}

	reader := bufio.NewReader(os.Stdin)

	output := "\n" +
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
	commitHash(artifactHash, owner, string(passphrase), filename)
	fmt.Println("")
	fmt.Println("Artifact:\t\t", filename)
	fmt.Println("Hash:\t\t", artifactHash)
	fmt.Println("Date:\t\t", time.Now())
	fmt.Println("Signer:\t\t", owner)
}

// VerifyAll => main entry point from cli
// unwraps potential list input for verify()
func VerifyAll(files []string) {
	for _, file := range files {
		verify(file)
	}
}
func verify(filename string) {

	var artifactHash string

	// TODO: make this switch available for all functions
	if strings.HasPrefix(filename, "docker:") {

		artifactHash = getDockerHash(filename)
		//fmt.Printf("Docker: Not yet implemented\n")
		//os.Exit(1)
	} else {
		artifactHash = strings.TrimSpace(hash(filename))
	}

	verified, owner, timestamp := verifyHash(artifactHash)
	fmt.Println("File:\t", filename)
	fmt.Println("Hash:\t", artifactHash)
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
