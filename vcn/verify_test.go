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
	"os/exec"
	"testing"
)

func TestVerifyDocker(t *testing.T) {
	// make sure the docker verification works e2e
	const helloWorld string = "hello-world"

	// Run the Docker hello world example
	// this ensures that the image will be present
	exec.Command("docker", "run", helloWorld)
	//fmt.Println(c)

	h := getDockerHash(fmt.Sprintf("docker:%s", helloWorld))

	verified, _, _ := verifyHash(h)

	if verified != true {
		t.Error(fmt.Sprintf("Could not verify: <docker:%s>", helloWorld))
	}

}
