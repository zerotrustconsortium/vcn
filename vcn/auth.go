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
	"strings"
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

type PublisherExistsResponse struct {
	Exists bool `json:"exists"`
}
type PublisherExistsParams struct {
	Email string `url:"email"`
}
type PublisherResponse struct {
	Authorities []string `json:"authorities"`
	Email       string   `json:"email"`
	FirstName   string   `json:"firstName"`
	LastName    string   `json:"lastName"`
}

func CheckPublisherExists(email string) (ret bool) {

	email = strings.TrimSpace(email)

	params := &PublisherExistsParams{Email: email}
	response := new(PublisherExistsResponse)
	restError := new(Error)

	r, err := sling.New().
		Get(PublisherEndpoint()+"/exists").
		QueryStruct(params).
		Receive(&response, restError)

	if err != nil {
		fmt.Printf(err.Error())
		return false
	}
	if r.StatusCode != 200 {

		fmt.Printf(fmt.Sprintf("request failed: %s (%d)",
			restError.Message, restError.Status))
		return false
	}

	return response.Exists
}

func CheckToken(token string) (ret bool) {

	if token == "" {
		return false
	}
	// TODO: change api call to real check token
	_, status := CheckPublisherIsVerified(token)

	//fmt.Println(status)
	if status == 0 {
		return true
	}
	return false

}

func Authenticate(email string, password string) (ret bool, code int) {

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
		// TODO: DEBUG LOG LEVEL
		fmt.Printf("request failed: %s (%d)\n", authError.Message,
			authError.Status)
		return false, authError.Status

	}
	err = ioutil.WriteFile(TokenFile(), []byte(token.Token),
		os.FileMode(0600))
	if err != nil {
		log.Fatal(err)
	}

	return true, 0

}

func CheckPublisherIsVerified(token string) (ret bool, status int) {

	restError := new(Error)
	response := new(PublisherResponse)

	r, err := sling.New().
		Get(PublisherEndpoint()).
		Add("Authorization", "Bearer "+token).
		Receive(&response, restError)

	if err != nil {
		// TODO DEBUG LEVEL
		//fmt.Printf(err.Error())
		return false, 500
	}
	if r.StatusCode != 200 {

		fmt.Printf(fmt.Sprintf("request failed: %s (%d)",
			restError.Message, restError.Status))
		return false, restError.Status
	}

	for _, el := range response.Authorities {
		if el == ROLE_CONFIRMED_USER() {
			return true, 0
		}
	}

	return false, 404
}

// Register creates an Account with vChain.us
func Register(email string, accountPassword string) (ret bool, code int) {

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
		// TODO debug log
		log.Printf("request failed: %s (%d)", authError.Message, authError.Status)

		return false, authError.Status
	}
	return true, 0
}

func WaitForConfirmation(email string, password string, maxRounds uint64,
	pollInterval time.Duration) error {

	fmt.Println("hier!")
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
