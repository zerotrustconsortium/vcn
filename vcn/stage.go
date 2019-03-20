/*
 * Copyright (c) 2018-2019 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package main

import (
	"log"
	"os"
)

type Stage int64

const (
	StageProduction Stage = 0
	StageStaging    Stage = 1
	StageTest       Stage = 2
)

func StageEnvironment() Stage {
	switch os.Getenv("STAGE") {
	case "STAGING":
		return StageStaging
	case "TEST":
		return StageTest
	case "PRODUCTION":
		return StageProduction
	default:
		return StageProduction
	}
}

func StageName(stage Stage) (name string) {
	switch stage {
	case StageProduction:
		return "PRODUCTION"
	case StageStaging:
		return "STAGING"
	case StageTest:
		return "TEST"
	default:
		log.Fatal("unsupported stage", name)
		return ""
	}
}
