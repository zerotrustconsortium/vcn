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

func APIEndpoint(method string) string {
	return "https://api.vchain.us/v1/" + method
}

func MainNetEndpoint() string {
	return "https://main.vchain.us"
}

func ProofContractAddress() string {
	return "0x66ccf074254cb0eb8d9e8020d8e777406a1d9cbb"
}
