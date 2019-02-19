/*
 * Copyright (c) 2018-2019 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package main

import (
	"fmt"
	"log"

	"github.com/dghubble/sling"
)

type ArtifactRequest struct {
	Name string `json:"name"`
	Hash string `json:"hash"`
}

type PagedArtifactResponse struct {
	Content []ArtifactResponse `json:"content"`
}

type ArtifactResponse struct {
	Name            string   `json:"name"`
	Hash            string   `json:"hash"`
	Level           int      `json:"level"`
	Visibility      string   `json:"visibility"`
	Status          string   `json:"status"`
	IntegrityChecks []string `json:"integrityChecks"`
}

func (a ArtifactResponse) String() string {

	return fmt.Sprintf("File:\t%s\nHash:\t%s\nStatus:\t%s\n\n",
		a.Name, a.Hash, a.Status)
}

func CreateArtifact(walletAddress string, name string, hash string) error {
	restError := new(Error)
	token, err := LoadToken()
	if err != nil {
		log.Fatal(err)
	}
	r, err := sling.New().
		Post(ArtifactEndpoint(walletAddress)).
		Add("Authorization", "Bearer "+token).
		BodyJSON(ArtifactRequest{
			Name: name,
			Hash: hash,
		}).Receive(nil, restError)
	if err != nil {
		return err
	}
	if r.StatusCode != 200 {
		return fmt.Errorf("request failed: %s (%d)", restError.Message,
			restError.Status)
	}
	return nil
}

func LoadArtifactsForCurrentWallet() ([]ArtifactResponse, error) {
	publicKey, err := PublicKeyForLocalWallet()
	if err != nil {
		return nil, err
	}
	return LoadArtifacts(publicKey)
}

func LoadArtifacts(walletAddress string) ([]ArtifactResponse, error) {
	response := new(PagedArtifactResponse)
	restError := new(Error)
	token, err := LoadToken()
	if err != nil {
		log.Fatal(err)
	}
	r, err := sling.New().
		Get(ArtifactEndpoint(walletAddress)).
		Add("Authorization", "Bearer "+token).
		Receive(&response, restError)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, fmt.Errorf("request failed: %s (%d)",
			restError.Message, restError.Status)
	}
	return response.Content, nil
}
