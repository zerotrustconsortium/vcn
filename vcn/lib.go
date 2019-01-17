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
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"log"
)

func hash(filename string) string {

	// TODO: large files
	// https://stackoverflow.com/questions/15879136/how-to-calculate-sha256-file-checksum-in-go
	hasher := sha256.New()
	s, err := ioutil.ReadFile(filename)
	hasher.Write(s)
	if err != nil {
		log.Fatal(err)
	}

	hash := hex.EncodeToString(hasher.Sum(nil))
	//fmt.Println("Hash: ", hash)

	return hash
}
