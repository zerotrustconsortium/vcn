/*
 * Copyright (c) 2018 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 * Built on top of CLI (MIT license)
 * https://github.com/urfave/cli#overview
 */

package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "vcn"
	app.Version = "0.0.1"

	app.Commands = []cli.Command{
		{
			Category: "Artifact actions",
			Name:     "verify",
			Aliases:  []string{"v"},
			Usage:    "verify against blockchain",
			Action: func(c *cli.Context) error {
				fmt.Println("verified artifact: ", c.Args().First())
				verify(c.Args().First())
				return nil
			},
			Flags: []cli.Flag{
				cli.BoolFlag{Name: "hash"},
			},
		},
		{
			Category: "Artifact actions",
			Name:     "commit",
			Aliases:  []string{"c"},
			Usage:    "commit in blockchain",
			Action: func(c *cli.Context) error {
				fmt.Println("committed artifact: ", c.Args().First())
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

type response struct {
	Message string `json:"message"`
}

func verify(filename string) {
	fmt.Println("File: ", filename)

	// TODO: large files
	// https://stackoverflow.com/questions/15879136/how-to-calculate-sha256-file-checksum-in-go
	hasher := sha256.New()
	s, err := ioutil.ReadFile(filename)
	hasher.Write(s)
	if err != nil {
		log.Fatal(err)
	}

	hash := hex.EncodeToString(hasher.Sum(nil))
	fmt.Println("Hash: ", hash)

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
	fmt.Println("Trust status:", verification.Message)

}
