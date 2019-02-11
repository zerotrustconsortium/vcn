package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"syscall"

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
