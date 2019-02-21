/*
 * Copyright (c) 2018-2019 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 */

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/dghubble/sling"
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

func Authenticate(email string, password string) (ret bool) {

	token := new(TokenResponse)
	authError := new(Error)

	r, err := sling.New().
		Post(PublisherEndpoint()+"/auth").
		BodyJSON(AuthRequest{Email: email, Password: password}).
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

	return true

}

// Register creates an Account with vChain.us
func Register(email string, accountPassword string) {

	authError := new(Error)
	//var apiError string

	r, err := sling.New().
		Post(PublisherEndpoint()).
		BodyJSON(AuthRequest{Email: email, Password: accountPassword}).
		Receive(nil, authError)
	if err != nil {
		log.Fatal(err)
	}
	if r.StatusCode != 200 {
		//GET-v1-artifact-404

		log.Printf("request failed: %s (%d)", authError.Message, authError.Status)
		PrintErrorURLByEndpoint(PublisherEndpoint(), "POST", authError.Status)
		os.Exit(1)
	}
}

func WaitForConfirmation(email string, password string, maxRounds uint64,
	pollInterval time.Duration) error {
	token := new(TokenResponse)
	authError := new(Error)
	for i := uint64(0); i < maxRounds; i++ {
		r, err := sling.New().
			Post(PublisherEndpoint()+"/auth").
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
