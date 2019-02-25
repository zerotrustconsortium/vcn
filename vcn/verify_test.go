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
	"testing"
)

func TestVerifyDocker(t *testing.T) {
	// make sure the docker verification works e2e
	const helloWorld string = "hello-world"

	// Run the Docker hello world example
	// this ensures that the image will be present
	// WON'T WORK INSIDE DOCKERIZED SETUP ANYHOW
	//exec.Command("docker", "run", helloWorld)
	//h := getDockerHash(fmt.Sprintf("docker:%s", helloWorld))
	// simply take the known hash of the <hello-world> image for now; no meaningful test actually
	h := "fce289e99eb9bca977dae136fbe2a82b6b7d4c372474c9235adc1741675f587e"

	verified, _, _ := verifyHash(h)

	if verified != true {
		t.Error(fmt.Sprintf("Could not verify: <docker:%s>", helloWorld))
	}

}
