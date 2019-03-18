package main

import (
	"log"
	"sync"
)

const VcnVersion = "0.3.3"
const VcnDirectoryPermissions = 0700

var WG sync.WaitGroup

type Level int64
type Status int64
type Visibility int64

const (
	LevelDisabled         Level = -1
	LevelUnknown          Level = 0
	LevelEmailVerified    Level = 1
	LevelSocialVerified   Level = 2
	LevelIdVerified       Level = 3
	LevelLocationVerified Level = 4
	LevelVchain           Level = 99
)

const (
	StatusTrusted     Status = 0
	StatusUntrusted   Status = 1
	StatusUnknown     Status = 2
	StatusUnsupported Status = 3
)

const (
	VisibilityPublic  Visibility = 0
	VisibilityPrivate Visibility = 1
)

func LevelName(level Level) (name string) {
	switch level {
	case LevelDisabled:
		return "DISABLED"
	case LevelUnknown:
		return "UNKNOWN"
	case LevelEmailVerified:
		return "EMAIL_VERIFIED"
	case LevelSocialVerified:
		return "SOCIAL_VERIFIED"
	case LevelIdVerified:
		return "ID_VERIFIED"
	case LevelLocationVerified:
		return "LOCATION_VERIFIED"
	case LevelVchain:
		return "VCHAIN"
	default:
		log.Fatal("unsupported level", name)
		return "";
	}
}

func StatusName(status Status) (name string) {
	switch status {
	case StatusTrusted:
		return "TRUSTED"
	case StatusUntrusted:
		return "UNTRUSTED"
	case StatusUnknown:
		return "UNKNOWN"
	case StatusUnsupported:
		return "UNSUPPORTED"
	default:
		log.Fatal("unsupported status", name)
		return "";
	}
}

func VisibilityName(visibility Visibility) (name string) {
	switch visibility {
	case VisibilityPublic:
		return "PUBLIC"
	case VisibilityPrivate:
		return "PRIVATE"
	default:
		log.Fatal("unsupported visibility", name)
		return "";
	}
}

func VisibilityForFlag(public bool) (visibility Visibility) {
	if public {
		return VisibilityPublic
	} else {
		return VisibilityPrivate
	}
}
