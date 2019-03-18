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
	"io"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"bou.ke/monkey"
)

// of course this is a temporary measure;-)
const USER = "mathias@vchain.us"
const PASSWORD = "***"
const PASSPHRASE = "WHATEVER"

func IgnoreTestLoginVerifiedUser(t *testing.T) {

	// fmt.Println(VcnDirectory())

	tdir, err := ioutil.TempDir("", "vcn-testing")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tdir) // clean up

	monkey.Patch(
		VcnDirectory,
		func() string {
			return tdir
		},
	)

	CreateVcnDirectories()

	// fmt.Println(VcnDirectory())

	in, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer in.Close()

	_, err = io.WriteString(
		in,
		fmt.Sprintf("%s\n%s\n%s\n%s\n",
			USER, PASSWORD, PASSPHRASE, PASSPHRASE))
	if err != nil {
		t.Fatal(err)
	}

	_, err = in.Seek(0, os.SEEK_SET)
	if err != nil {
		t.Fatal(err)
	}

	login(in)
	// TODO: parse STDOUT

}
func TestVerifyOutput(t *testing.T) {

	/*
		expectedURL := "htertps://github.com/vchain-us/vcn/wiki/Errors#publisher-post-412"

		res := PublisherEndpoint()
		verb := "pOsT" // should do lowercase
		status := 412

		actualURL := formatErrorURLByEndpoint(res, verb, status)

		if expectedURL != actualURL {
			t.Error(fmt.Sprintf("formatErrorURLByEndpoint() does not match [%s != %s]", expectedURL, actualURL))
		}
	*/

}
