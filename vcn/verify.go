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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/fatih/color"
)

type response struct {
	Message string `json:"message"`
}

func verify(filename string) {
	fmt.Println("File: ", filename)

	hash := hash(filename)

	url := "http://api.vchain.us/v1/files/" + hash

	vcnClient := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, getErr := vcnClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	verification := response{}
	jsonErr := json.Unmarshal(body, &verification)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	// TO DO: Trust level= ?? mb, also this should also go to STDOUT,
	// and also handle STDERR so it can be fully scripted

	trust := strings.TrimSpace(verification.Message)

	color.Set(color.FgHiWhite, color.BgCyan, color.Bold)
	fmt.Print("Trust status:", trust)
	color.Unset()

	fmt.Println()
}
