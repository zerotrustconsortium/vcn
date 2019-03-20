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

	"github.com/sirupsen/logrus"
)

type VerifyArtifactTrackingEventRequest struct {
	Client   string `json:"client"`
	Filename string `json:"filename"`
	Hash     string `json:"hash"`
	Url      string `json:"url"`
}

type SignArtifactTrackingEventRequest struct {
	Filename string `json:"filename"`
	Hash     string `json:"hash"`
	Name     string `json:"name"`
	Url      string `json:"url"`
}

type PublisherTrackingEventRequest struct {
	Name string `json:"name"`
}

func TrackVerify(hash string, filename string) (err error) {
	LOG.WithFields(logrus.Fields{
		"hash":     hash,
		"filename": filename,
	}).Trace("TrackVerify")
	restError := new(Error)
	token, _ := LoadToken()
	r, err := NewSling(token).
		Post(TrackingEvent() + "/verify").
		BodyJSON(VerifyArtifactTrackingEventRequest{
			Client:   VcnClientName(),
			Filename: filename,
			Hash:     hash,
		}).Receive(nil, &restError)
	if err != nil {
		return err
	}
	if r.StatusCode != 200 {
		return fmt.Errorf("TrackVerify failed: %+v", restError)
	}
	return nil
}

func TrackPublisher(event string) (err error) {
	LOG.WithFields(logrus.Fields{
		"event": event,
	}).Trace("TrackPublisher")
	restError := new(Error)
	token, err := LoadToken()
	if err != nil {
		return err
	}
	r, err := NewSling(token).
		Post(TrackingEvent() + "/publisher").
		BodyJSON(PublisherTrackingEventRequest{
			Name: event,
		}).Receive(nil, &restError)
	if err != nil {
		return err
	}
	if r.StatusCode != 200 {
		return fmt.Errorf("TrackPublisher failed: %+v", restError)
	}
	return nil
}

func TrackSign(hash string, filename string, status Status) (err error) {
	LOG.WithFields(logrus.Fields{
		"hash":     hash,
		"filename": filename,
		"status":   status,
	}).Trace("TrackSign")
	restError := new(Error)
	token, err := LoadToken()
	if err != nil {
		return err
	}
	r, err := NewSling(token).
		Post(TrackingEvent() + "/sign").
		BodyJSON(SignArtifactTrackingEventRequest{
			Name:     StatusName(status),
			Hash:     hash,
			Filename: filename,
		}).Receive(nil, &restError)
	if err != nil {
		return err
	}
	if r.StatusCode != 200 {
		return fmt.Errorf("TrackSign failed: %+v", restError)
	}
	return nil
}
