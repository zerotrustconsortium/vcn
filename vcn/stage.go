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
