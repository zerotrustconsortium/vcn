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
	"math/big"
	"os"
	"time"

	"github.com/fatih/color"
)

type Level int

const (
	DISABLED          Level = -1
	UNKNOWN           Level = 0
	EMAIL_VERIFIED    Level = 1
	SOCIAL_VERIFIED   Level = 2
	ID_VERIFIED       Level = 3
	LOCATION_VERIFIED Level = 4
	VCHAIN            Level = 99
)

type Status int

const (
	OK          Status = 0
	UNSUPPORTED Status = 1
	UNTRUSTED   Status = 2
)

func StyleAffordance() (color.Attribute, color.Attribute, color.Attribute) {
	return color.FgHiBlue, color.BgWhite, color.Bold
}
func StyleError() (color.Attribute, color.Attribute, color.Attribute) {
	return color.FgHiRed, color.BgHiWhite, color.Bold
}
func ErrorWikiURL() string {
	return "https://github.com/vchain-us/vcn/wiki/Errors#"
}
func DashboardURL() string {
	return "https://dashboard.staging.vchain.us"
}
func MainNetEndpoint() string {
	return "https://main.staging.vchain.us"
}

func FoundationEndpoint() string {
	return "https://api.staging.vchain.us/foundation"
}

func TrackingEvent() string {
	return FoundationEndpoint() + "/v1/tracking-event"
}

func TokenCheckEndpoint() string {
	return PublisherEndpoint() + "/auth/check"
}

func PublisherEndpoint() string {
	return FoundationEndpoint() + "/v1/publisher"
}

func ROLE_CONFIRMED_USER() string {
	return "ROLE_CONFIRMED_USER"
}

func WalletEndpoint() string {
	return FoundationEndpoint() + "/v1/wallet"
}

func ArtifactEndpoint(walletAddress string) string {
	return FoundationEndpoint() + "/v1/artifact?wallet-address=" + walletAddress
}

func ProofContractAddress() string {
	return "0x66ccf074254cb0eb8d9e8020d8e777406a1d9cbb"
}

func AssetsRelayContractAddres() string {
	return "0xf1d4b9fe8290bb5718db5d46c313e7b266570c21"
}

func TxVerificationRounds() uint64 {
	return 10
}

func PollInterval() time.Duration {
	return 1 * time.Second
}

func VcnDirectory() string {
	return os.Getenv("HOME") + "/.vcn"
}

func WalletDirectory() string {
	return VcnDirectory() + "/wallets"
}

func TokenFile() string {
	return VcnDirectory() + "/t"
}

func GasPrice() *big.Int {
	return big.NewInt(0)
}

func GasLimit() uint64 {
	return 20000000
}
