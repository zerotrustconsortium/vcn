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

func TestGetStatusName(t *testing.T) {

	name := getStatusName(int(UNSUPPORTED))
	name += getStatusName(int(OK))
	name += getStatusName(int(UNTRUSTED))

	if name != "UNSUPPORTEDOKUNTRUSTED" {
		t.Error("Status enumeration seems broken")
	}
}
func TestErrorURLComposition(t *testing.T) {

	expectedURL := "https://github.com/vchain-us/vcn/wiki/Errors#publisher-post-412"

	res := PublisherEndpoint()
	verb := "pOsT" // should do lowercase
	status := 412

	actualURL := formatErrorURLByEndpoint(res, verb, status)

	if expectedURL != actualURL {
		t.Error(fmt.Sprintf("formatErrorURLByEndpoint() does not match [%s != %s]", expectedURL, actualURL))
	}

}
