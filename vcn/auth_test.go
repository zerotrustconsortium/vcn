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

func TestCheckTokenNoInput(t *testing.T) {

	token := ""
	ret := CheckToken(token)

	if ret != false {
		t.Error(fmt.Sprintf(`CheckToken() with empty string input must return false`))
	}

}
