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
	"os"

	"github.com/dghubble/sling"
)

type ArtifactRequest struct {
	Name       string `json:"name"`
	Hash       string `json:"hash"`
	Filename   string `json:"filename"`
	Url        string `json:"url"`
	License    string `json:"license"`
	Visibility string `json:"visibility"`
	Status     string `json:"status"`
	MetaHash   string `json:"metaHash"`
}

type PagedArtifactResponse struct {
	Content []ArtifactResponse `json:"content"`
}

type ArtifactResponse struct {
	Name       string `json:"name"`
	Hash       string `json:"hash"`
	Filename   string `json:"filename"`
	Level      int    `json:"level"`
	Visibility string `json:"visibility"`
	Status     string `json:"status"`
	Publisher  string `json:"publisher"`
	CreatedAt  string `json:"createdAt"`
}

func (a ArtifactResponse) String() string {

	return fmt.Sprintf("File:\t%s\nHash:\t%s\nStatus:\t%s\n\n",
		a.Name, a.Hash, a.Status)
}

func CreateArtifact(walletAddress string, name string, hash string, visibility Visibility, status Status) error {
	restError := new(Error)
	token, err := LoadToken()
	if err != nil {
		fmt.Printf("\n%s\n", err.Error())
		PrintErrorURLCustom("sign", 404)
		os.Exit(1)
	}
	metaHash, err := hashAsset(hash)
	if err != nil {
		log.Fatal("unable to hash asset", err)
	}
	r, err := sling.New().
		Post(ArtifactEndpointForWallet(walletAddress)).
		Add("Authorization", "Bearer "+token).
		BodyJSON(ArtifactRequest{
			Name:       name,
			Hash:       hash,
			Filename:   name,
			Visibility: visibilityName(visibility),
			Status:     statusName(int(status)),
			MetaHash:   metaHash,
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
		Get(ArtifactEndpointForWallet(walletAddress)).
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

func LoadArtifactsForHash(hash string, metahash string) (*ArtifactResponse, error) {
	response := new(ArtifactResponse)
	restError := new(Error)
	token, err := LoadToken()
	if err != nil {
		log.Fatal(err)
	}
	r, err := sling.New().
		Get(ArtifactEndpoint() + "/" + hash + "/" + metahash).
		Add("Authorization", "Bearer "+token).
		Receive(&response, restError)
	if err != nil {
		return nil, err
	}
	if r.StatusCode == 404 {
		return nil, nil
	}
	if r.StatusCode != 200 {
		return nil, fmt.Errorf("request failed: %s (%d)",
			restError.Message, restError.Status)
	}
	return response, nil
}
