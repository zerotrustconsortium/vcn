package main

import (
	"log"
	"sync"
)

var VcnVersion = "0.3.3"
var WG sync.WaitGroup

type Level int
type Status int
type Visibility int

const (
	LEVEL_DISABLED          Level = -1
	LEVEL_UNKNOWN           Level = 0
	LEVEL_EMAIL_VERIFIED    Level = 1
	LEVEL_SOCIAL_VERIFIED   Level = 2
	LEVEL_ID_VERIFIED       Level = 3
	LEVEL_LOCATION_VERIFIED Level = 4
	LEVEL_VCHAIN            Level = 99
)

const (
	STATUS_TRUSTED     Status = 0
	STATUS_UNTRUSTED   Status = 1
	STATUS_UNKNOWN     Status = 2
	STATUS_UNSUPPORTED Status = 3
)

const (
	VISIBILITY_PUBLIC  Visibility = 0
	VISIBILITY_PRIVATE Visibility = 1
)

func levelName(level int) (name string) {
	switch level {
	case int(LEVEL_DISABLED):
		return "DISABLED"
	case int(LEVEL_UNKNOWN):
		return "UNKNOWN"
	case int(LEVEL_EMAIL_VERIFIED):
		return "EMAIL_VERIFIED"
	case int(LEVEL_SOCIAL_VERIFIED):
		return "SOCIAL_VERIFIED"
	case int(LEVEL_ID_VERIFIED):
		return "ID_VERIFIED"
	case int(LEVEL_LOCATION_VERIFIED):
		return "LOCATION_VERIFIED"
	case int(LEVEL_VCHAIN):
		return "VCHAIN"
	default:
		log.Fatal("unsupported level", name)
		return "";
	}
}

func statusName(status int) (name string) {
	switch status {
	case int(STATUS_TRUSTED):
		return "TRUSTED"
	case int(STATUS_UNTRUSTED):
		return "UNTRUSTED"
	case int(STATUS_UNKNOWN):
		return "UNKNOWN"
	case int(STATUS_UNSUPPORTED):
		return "UNSUPPORTED"
	default:
		log.Fatal("unsupported status", name)
		return "";
	}
}

func visibilityName(visibility Visibility) (name string) {
	switch visibility {
	case VISIBILITY_PUBLIC:
		return "PUBLIC"
	case VISIBILITY_PRIVATE:
		return "PRIVATE"
	default:
		log.Fatal("unsupported visibility", name)
		return "";
	}
}

func visibilityForFlag(public bool) (visibility Visibility) {
	if public {
		return VISIBILITY_PUBLIC
	} else {
		return VISIBILITY_PRIVATE
	}
}
