/*
 * Copyright (c) 2018-2019 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/fatih/color"
	"github.com/yalp/jsonpath"
	"golang.org/x/crypto/ssh/terminal"
)

func firstFile(dir string) (io.Reader, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		return os.Open(dir + "/" + f.Name())
	}
	return nil, fmt.Errorf("empty directory: %s", dir)
}

func contains(xs []string, x string) bool {
	for _, a := range xs {
		if a == x {
			return true
		}
	}
	return false
}

func readPassword(msg string) (string, error) {
	fmt.Print(msg)
	password, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Println(".")
	if err != nil {
		return "", err
	}
	return string(password), nil
}

func CreateVcnDirectories() {
	if err := os.MkdirAll(WalletDirectory(),
		os.FileMode(VcnDirectoryPermissions)); err != nil {
		log.Fatal(err)
	}
}

// PrintErrorURLCustom takes custom domain and status code
func PrintErrorURLCustom(domain string, code int) {

	fmt.Print("Get help for this error at:\n")

	color.Set(StyleError())
	fmt.Print(formatErrorURLCustom(domain, code))
	color.Unset()

	fmt.Println()
	return

}
func formatErrorURLCustom(domain string, status int) string {

	errorPage := ErrorWikiURL()

	return fmt.Sprintf("%s%s-%d", errorPage, domain, status)

}

// PrintErrorURLByEndpoint takes API errors and creates github wiki links
func PrintErrorURLByEndpoint(resource string, verb string, status int) {

	fmt.Print("Get help for this error at:\n")

	color.Set(StyleError())
	fmt.Print(formatErrorURLByEndpoint(resource, verb, status))
	color.Unset()

	fmt.Println()
	return

}
func formatErrorURLByEndpoint(resource string, verb string, status int) string {

	errorPage := ErrorWikiURL()

	// get last part of endpoint
	x := strings.Split(resource, "/")
	resource = x[len(x)-1]

	return fmt.Sprintf("%s%s-%s-%d", errorPage, resource, strings.ToLower(verb), status)

}

func getDockerHash(param string) (hash string) {

	dockerID := strings.Replace(param, "docker:", "", 1)

	// TODO: sanitize even further
	// so far, let's check dockerID is a string without whitespaces
	dockerID = strings.Replace(dockerID, " ", "", -1)

	/*

		hash = string(output)
		hash = strings.Replace(hash, `"`, ``, -1)
		hash = strings.Replace(hash, "sha256:", "", 1)
	*/

	cmd := exec.Command("docker", "inspect", dockerID)
	cmdOutput, err := cmd.Output()
	if err != nil {
		fmt.Printf(fmt.Sprintf("Failed to execute docekr inspect command."))
		fmt.Printf(err.Error())
		PrintErrorURLCustom("docker", 500)
		os.Exit(1)
	}

	// var dockerInspect interface {}
	dockerIDFilter, err := jsonpath.Prepare("$..Id")
	if err != nil {
		panic(err)
	}
	var data interface{}
	if err = json.Unmarshal(cmdOutput, &data); err != nil {
		panic(err)
	}
	out, err := dockerIDFilter(data)
	if err != nil {
		panic(err)
	}

	// out is an interface which needs to be coreced into string array before
	dockerHash := out.([]interface{})[0]
	dockerHashStr := strings.TrimSpace(strings.Replace(fmt.Sprint(dockerHash), "sha256:", "", 1))

	return dockerHashStr
}

func hashAsset(assetHash string) (metadataHash string, err error) {
	verification, err := BlockChainVerify(assetHash)
	if err != nil {
		return "", err
	}
	metadata := fmt.Sprintf("%s-%d-%d-%d",
		verification.Owner.Hex(),
		int64(verification.Level),
		int64(verification.Status),
		int64(verification.Timestamp.Unix()))
	metadataHashAsBytes := sha256.Sum256([]byte(metadata))
	return fmt.Sprintf("%x", metadataHashAsBytes), nil
}
