package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"syscall"

	"github.com/fatih/color"
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

func createDirectoryInfrastructure() {
	err := os.MkdirAll(WalletDirectory(), os.FileMode(0700))
	if err != nil {
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
