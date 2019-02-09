package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"syscall"
	"time"

	"github.com/dghubble/sling"
	"golang.org/x/crypto/ssh/terminal"
)

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenResponse struct {
	Token string `token:"token"`
}

type Error struct {
	Status    int    `json:"status"`
	Message   string `json:"message"`
	Path      string `json:"path"`
	Timestamp string `json:"timestamp"`
	Error     string `json:"error"`
}

func Authenticate(email string) {
	token := new(TokenResponse)
	authError := new(Error)
	fmt.Print("Password:")
	password, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Println(".")
	if err != nil {
		log.Fatal(err)
	}
	r, err := sling.New().
		Post(PublisherEndpoint() + "/auth").
		BodyJSON(AuthRequest{Email: email, Password: string(password)}).
		Receive(token, authError)
	if err != nil {
		log.Fatal(err)
	}
	if r.StatusCode != 200 {
		log.Fatalf("request failed: %s (%d)", authError.Message,
			authError.Status)
	}
	err = ioutil.WriteFile(TokenFile(), []byte(token.Token),
		os.FileMode(0600))
	if err != nil {
		log.Fatal(err)
	}
}

func Register(email string) {
	authError := new(Error)
	fmt.Print("Password:")
	password, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Println(".")
	if err != nil {
		log.Fatal(err)
	}
	r, err := sling.New().
		Post(PublisherEndpoint()).
		BodyJSON(AuthRequest{Email: email, Password: string(password)}).
		Receive(nil, authError)
	if err != nil {
		log.Fatal(err)
	}
	if r.StatusCode != 200 {
		log.Fatalf("request failed: %s (%d)", authError.Message,
			authError.Status)
	}
	log.Println("We've sent you an email to: ", email,
		"Click the link and you will be automatically logged in")
	err = WaitForConfirmation(email, string(password),
		60, 2*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	// @ameingast: check for existing key-stores
	CreateKeystore(string(password))
	SyncKeys()
}

func WaitForConfirmation(email string, password string, maxRounds uint64,
	pollInterval time.Duration) error {
	token := new(TokenResponse)
	authError := new(Error)
	for i := uint64(0); i < maxRounds; i++ {
		r, err := sling.New().
			Post(PublisherEndpoint() + "/auth").
			BodyJSON(AuthRequest{Email: email, Password: password}).
			Receive(token, authError)
		if err != nil {
			return err
		}
		if r.StatusCode == 400 {
			time.Sleep(pollInterval)
			continue
		}
		if r.StatusCode == 200 {
			err = ioutil.WriteFile(TokenFile(), []byte(token.Token),
				os.FileMode(0600))
			if err != nil {
				return err
			}
			return nil
		}
		return fmt.Errorf("wait for confirmation failed: %s (%d)",
			authError.Message, authError.Status)
	}
	return fmt.Errorf("confirmation timed out")
}

func LoadToken() (jwtToken string, err error) {
	contents, err := ioutil.ReadFile(TokenFile())
	if err != nil {
		return "", err
	}
	return string(contents), nil
}
