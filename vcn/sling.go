package main

import (
	"github.com/dghubble/sling"
)

func NewSling(token string) (s *sling.Sling) {
	s = sling.New()
	if token != "" {
		if token != "" {
			s = s.Add("Authorization", "Bearer "+token)
		}
	}
	return s
}
