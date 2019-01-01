/*
 * Copyright (c) 2018-2019 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 * Built on top of CLI (MIT license)
 * https://github.com/urfave/cli#overview
 */

package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func commit(filename string, auth string) {

	type response struct {
		Message         string `json:"message"`
		TransactionHash string `json:"transaction-hash"`
		Block           string `json:"block"`
	}

	fmt.Println("File: ", filename)

	hash := hash(filename)
	b64data := []byte(auth)
	b64str := base64.StdEncoding.EncodeToString(b64data)

	url := "http://api.vchain.us/v1/files"

	var jsonStr = []byte(fmt.Sprintf(`{"hash":"%s", "owner":"na"}`, hash))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", b64str))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	//fmt.Println("response Status:", resp.Status)
	//fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)

	//fmt.Println("response Body:", string(body))

	commit := response{}
	jsonErr := json.Unmarshal(body, &commit)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	fmt.Println("Commit status:", commit.Message)

	fmt.Println("Block:", commit.Block)

}
