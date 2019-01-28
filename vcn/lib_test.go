/*
 * Copyright (c) 2018-2019 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package main

import (
	"testing"
)

const expectedHash string = "181210f8f9c779c26da1d9b2075bde0127302ee0e3fca38c9a83f5b1dd8e5d3b"

func TestHash(t *testing.T) {
	if hash("../resources/testHash.example") != expectedHash {
		t.Error(`hash("../resources/testHash.example") does not match`)
	}

}
