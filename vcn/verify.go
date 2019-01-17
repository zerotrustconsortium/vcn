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

func verifyAll(files []string) {
	// find . -type f -name "*.go" | xargs -I % vcn verify %
	//for --> verify
	//fmt.Println(files)
	for i := 0; i < len(files); i++ {
		//fmt.Println(i)
		verify(files[i])
	}
}

func verify(filename string) {

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

	fmt.Println("File:", filename)
	fmt.Println("Hash:", strings.TrimSpace(hash))
	fmt.Print("      ")
	if trust == "ok" {
		color.Set(color.FgHiWhite, color.BgCyan, color.Bold)
	} else if trust == "Not found" {
		color.Set(color.FgHiWhite, color.BgRed, color.Bold)
	}

	fmt.Print("Trust:", trust)
	color.Unset()

	fmt.Println()

	// how to use stderr and exit codes
	//fmt.Fprintln(os.Stderr, "err")
	//os.Exit(1)

}
