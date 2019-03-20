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
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/fatih/color"
	"github.com/pkg/browser"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"
)

var displayProgress = true

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
		var email string

		fmt.Print("Email address: ")
		_, err := fmt.Scanln(&email)
		if err != nil {
			log.Fatal(err)
		}
		email = strings.Trim(email, "\n")
		email = strings.Trim(email, "\r")

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

	_ = TrackPublisher("VCN_LOGIN")

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

		fmt.Println("Keystore successfully created. We are updating your user profile.\n" +
			"You will be able to sign your first asset in one minute")
		fmt.Println("Public key:\t", pubKey)
		fmt.Println("Keystore:\t", wallet)

	}

	//
	SyncKeys()

	fmt.Println("Login successful.")

	WG.Wait()

}

// Commit => "sign"
func Sign(filename string, state Status, visibility Visibility, quit bool) {

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
	var fileSize int64 = 0

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
		fi, err := os.Stat(filename);
		if err != nil {
			log.Fatal(err)
		}
		fileSize = fi.Size()
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

	if strings.ToLower(strings.TrimSpace(question)) != "y" {

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

	_ = TrackPublisher("VCN_SIGN")
	_ = TrackSign(artifactHash, filepath.Base(filename), state)

	// TODO: return and display: block #, trx #
	_, _ = commitHash(artifactHash, string(passphrase), filepath.Base(filename), fileSize, state, visibility)
	fmt.Println("")
	fmt.Println("Asset:\t", filename)
	fmt.Println("Hash:\t", artifactHash)
	// fmt.Println("Date:\t\t", time.Now())
	// fmt.Println("Signer:\t", "<pubKey>")

	WG.Wait()
	displayProgress = false
	if !quit {
		if _, err := fmt.Scanln(); err != nil {
			log.Fatal(err)
		}
	}
}

func VerifyAll(files []string, quit bool) {
	_ = TrackPublisher("VCN_VERIFY")
	var success = true
	for _, file := range files {
		success = success && verify(file)
	}
	if !quit {
		if _, err := fmt.Scanln(); err != nil {
			log.Fatal(err)
		}
	}
	if success {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}

func verify(filename string) (success bool) {
	var artifactHash string

	// TODO: make this switch available for all functions
	if strings.HasPrefix(filename, "docker:") {

		artifactHash = getDockerHash(filename)
		// fmt.Printf("Docker: Not yet implemented\n")
		// os.Exit(1)
	} else {
		artifactHash = strings.TrimSpace(hash(filename))
	}
	if err := TrackVerify(artifactHash, filename); err != nil {
		log.Fatal("TrackVerify failed", err)
	}
	verification, err := BlockChainVerify(artifactHash)
	if err != nil {
		log.Fatal("unable to verify hash", err)
	}
	fmt.Println("Asset:\t", filename)
	fmt.Println("Hash:\t", artifactHash)

	if verification.Owner != common.BigToAddress(big.NewInt(0)) {
		fmt.Println("Date:\t", verification.Timestamp)
		metaHash, err := hashAsset(artifactHash)
		if err != nil {
			log.Fatal("Unable to calculate metahash")
		}
		artifact, err := LoadArtifactsForHash(artifactHash, metaHash)
		if err != nil {
			log.Fatal("Unable to resolve metahash")
		}
		if artifact != nil && artifact.Visibility == "PUBLIC" {
			fmt.Println("Signer:\t", artifact.Publisher)
			fmt.Println("Name:\t", artifact.Name)
			fmt.Println("File:\t", artifact.Filename)
		} else {
			fmt.Println("Signer:\t", verification.Owner.Hex())
		}
		fmt.Println("Level:\t", LevelName(verification.Level))

	} else {
		fmt.Println("Signer:\t NA")
		fmt.Println("Level:\t NA")
	}
	fmt.Print("Status:\t ")
	if verification.Status == StatusTrusted {
		color.Set(StyleSuccess())
		success = true
	} else {
		color.Set(StyleError())
		success = false
	}
	fmt.Print(StatusName(verification.Status))
	color.Unset()
	fmt.Println()
	return success
}

func displayLatency() {
	i := 0
	for displayProgress {
		i++
		fmt.Printf("\033[2K\rIn progress %02dsec", i)
		// fmt.Println(i)
		time.Sleep(1 * time.Second)
	}
}
