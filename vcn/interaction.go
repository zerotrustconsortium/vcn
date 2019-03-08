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
	"github.com/sirupsen/logrus"
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

func login(in *os.File) {
	if in == nil {
		in = os.Stdin
	}
	reader := bufio.NewReader(in)

	// file system: token exists && api: token is valid
	// no => enter email
	//        api: publisher exists
	//        yes => enter pw
	//               authenticate()
	//               fails => retry pw entry up to 3 times
	//        no  => hint at registration
	// filesystem: keystore exists
	// no => createkeystore
	// synckeys

	token, _ := LoadToken()
	tokenValid := CheckToken(token)

	if tokenValid == false {

		fmt.Print("Email address: ")
		email, _ := reader.ReadString('\n')
		email = strings.TrimSuffix(email, "\n")

		LOG.WithFields(logrus.Fields{
			"email": email,
		}).Trace("Email entered")

		publisherExists := CheckPublisherExists(email)

		if publisherExists {

			LOG.WithFields(logrus.Fields{
				"email": email,
			}).Debug("Publisher exists")

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

				// handle not-displayed password via STDIN
				// as well as file input for automated tests
				var passwordString string

				fmt.Printf("Password%s: ", attempt)
				// TODO: solution for reading from file inputs whose compilation does not fail on windows
				// if terminal.IsTerminal(syscall.Stdin) {
				password, _ := terminal.ReadPassword(int(syscall.Stdin))
				passwordString = string(password)
				fmt.Println("")
				/*} else {

					passwordString, _ = reader.ReadString('\n')
					passwordString = strings.TrimSuffix(passwordString, "\n")

				}*/

				returnCode := 0
				authenticated, returnCode = Authenticate(email, passwordString)

				if returnCode > 0 {

					if returnCode == 401 {
						fmt.Println("Please enter a correct password.")

					} else if returnCode == 400 {

						PrintErrorURLCustom("customer-verification", 412)

						LOG.WithFields(logrus.Fields{
							"code": returnCode,
						}).Fatal("API request failed: Email not confirmed.")
					}
				}

			}

		} else {
			fmt.Println("It looks like you have not yet registered.")
			color.Set(StyleAffordance())
			fmt.Printf("Please create an account first at %s", DashboardURL())
			color.Unset()
			fmt.Println()
			dashboard()
			os.Exit(1)
		}

	}

	// track successful login as early as possible
	// fire a go routine for the tracking that shall not delay the main user interaction
	WG.Add(1)
	go publisherEventTracker("VCN_LOGIN")

	hasKeystore, err := HasKeystore()
	if err != nil {
		LOG.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Could not access keystore directory")
	}
	if hasKeystore == false {

		fmt.Println("You have no keystore set up yet.")
		fmt.Println("<vcn> will now do this for you and upload the public key to the platform.")

		color.Set(StyleAffordance())
		fmt.Print("Attention: Please pick a strong passphrase. There is no recovery possible.")
		color.Unset()
		fmt.Println()

		var keystorePassphrase string
		var keystorePassphrase2 string

		match := false
		counter := 0
		for match == false {

			counter++

			if counter == 4 {
				fmt.Println("Too many attempts failed.")
				PrintErrorURLCustom("password", 404)
				os.Exit(1)

			}

			// TODO: solution for reading from file inputs whose compilation does not fail on windows
			// if terminal.IsTerminal(syscall.Stdin) {

			keystorePassphrase, _ = readPassword("Keystore passphrase: ")
			keystorePassphrase2, _ = readPassword("Keystore passphrase (reenter): ")
			fmt.Println("")
			/*} else {

				keystorePassphrase, _ = reader.ReadString('\n')
				keystorePassphrase = strings.TrimSuffix(keystorePassphrase, "\n")

				keystorePassphrase2, _ = reader.ReadString('\n')
				keystorePassphrase2 = strings.TrimSuffix(keystorePassphrase2, "\n")
			}*/

			if keystorePassphrase == "" {
				fmt.Println("Your passphrase must not be empty.")
			} else if keystorePassphrase != keystorePassphrase2 {
				fmt.Println("Your two inputs did not match. Please try again.")
			} else {
				match = true
			}

		}

		pubKey, wallet := CreateKeystore(keystorePassphrase)

		fmt.Println("Keystore successfully created")
		fmt.Println("Public key:\t", pubKey)
		fmt.Println("Keystore:\t", wallet)

	}

	//
	SyncKeys()

	fmt.Println("Login successful.")

	WG.Wait()

}

// Commit => "sign"
func Sign(filename string, state Status) {

	// check for token
	token, _ := LoadToken()
	if token == "" {
		fmt.Println("You need to be logged in to sign.")
		fmt.Println("Proceed by authenticating yourself using <vcn auth>")
		// PrintErrorURLCustom("token", 428)
		os.Exit(1)
	}

	// keystore
	hasKeystore, _ := HasKeystore()
	if hasKeystore == false {
		fmt.Printf("You need a keystore to sign.\n")
		fmt.Println("Proceed by authenticating yourself using <vcn auth>")
		// PrintErrorURLCustom("keystore", 428)
		os.Exit(1)
	}

	var artifactHash string

	// artifact types
	if strings.HasPrefix(filename, "docker:") {

		artifactHash = getDockerHash(filename)
		// fmt.Printf("Docker: Not yet implemented\n")
		// os.Exit(1)

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

	WG.Add(1)
	go publisherEventTracker("VCN_SIGN")
	WG.Add(1)
	go artifactCommitTracker(artifactHash, filename, state)

	// TODO: return and display: block #, trx #
	_, _ = commitHash(artifactHash, string(passphrase), filename, state)
	fmt.Println("")
	fmt.Println("Asset:\t", filename)
	fmt.Println("Hash:\t", artifactHash)
	// fmt.Println("Date:\t\t", time.Now())
	// fmt.Println("Signer:\t", "<pubKey>")

	WG.Wait()
}

// VerifyAll => main entry point from cli
// unwraps potential list input for verify()
func VerifyAll(files []string) {

	WG.Add(1)
	go publisherEventTracker("VCN_VERIFY")

	for _, file := range files {
		verify(file)
	}

	WG.Wait()
}
func verify(filename string) {

	var artifactHash string

	// TODO: make this switch available for all functions
	if strings.HasPrefix(filename, "docker:") {

		artifactHash = getDockerHash(filename)
		// fmt.Printf("Docker: Not yet implemented\n")
		// os.Exit(1)
	} else {
		artifactHash = strings.TrimSpace(hash(filename))
	}

	// fire a go routine for the tracking that shall not delay the main user interaction
	WG.Add(1)
	go artifactVerifyTracker(artifactHash, filename)

	verified, owner, level, status, timestamp := verifyHash(artifactHash)

	LOG.WithFields(logrus.Fields{
		"verified": verified,
	}).Debug("verifyHash()")

	fmt.Println("Asset:\t", filename)
	fmt.Println("Hash:\t", artifactHash)

	if timestamp != 0 {
		fmt.Println("Date:\t", time.Unix(timestamp, 0))
	}
	if owner != "" {
		fmt.Println("Signer:\t", owner)
	}

	fmt.Println("Level:\t", getLevelName(int(level)))

	fmt.Print("Status:\t ")
	if status == int64(OK) {
		color.Set(StyleSuccess())
	} else {
		color.Set(StyleError())
		defer os.Exit(1)
	}
	fmt.Print(getStatusName(int(status)))
	color.Unset()

	fmt.Println()

	// wait for the asset tracker to put data to analytics
	WG.Wait()
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
