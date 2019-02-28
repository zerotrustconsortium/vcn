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
	"fmt"
	"log"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var log2 = logrus.New()

func main() {

	ll := os.Getenv("LOG_LEVEL")
	switch ll {
	case "INFO":
		log2.SetLevel(logrus.InfoLevel)
	case "TRACE":
		log2.SetLevel(logrus.TraceLevel)
	default:
		log2.SetLevel(logrus.WarnLevel)

	}

	log2.WithFields(logrus.Fields{
		"animal": "walrus",
		"size":   10,
	}).Info("A group of walrus emerges from the ocean")

	log2.WithFields(logrus.Fields{
		"animal": "walrus",
		"size":   10,
	}).Warn("A group of walrus emerges from the ocean")
	log2.Trace("A group of walrus emerges from the ocean")

	app := cli.NewApp()
	app.Name = "vcn"
	app.Usage = "code signing made easy"
	app.Version = "0.2.0"

	app.Commands = []cli.Command{
		// possible commands:
		// trace <artifact>
		// list <pubkey>
		// search <block>
		// display validators

		{
			Category: "Artifact actions",
			Name:     "verify",
			Aliases:  []string{"v"},
			Usage:    "Verify digital artifact against blockchain",
			Action: func(c *cli.Context) error {
				if c.NArg() == 0 {
					return fmt.Errorf("assets required")
				}
				VerifyAll(c.Args())
				return nil
			},
			Flags: []cli.Flag{
				cli.BoolFlag{Name: "hash"},
			},
		},
		{
			Category: "Artifact actions",
			Name:     "sign",
			Aliases:  []string{"s"},
			Usage:    "Sign commits and artifact's hash onto the blockchain",
			Action: func(c *cli.Context) error {
				if c.NArg() == 0 {
					return fmt.Errorf("filename or type:reference required")
				}
				Sign(c.Args().First(), "me")
				return nil
			},
		},
		{
			Category: "Artifact actions",
			Name:     "list",
			Aliases:  []string{"l"},
			Usage:    "List your signed artifacts",
			Action: func(c *cli.Context) error {
				artifacts, err := LoadArtifactsForCurrentWallet()
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println("Artifacts:\n", artifacts)
				return nil
			},
		},
		{
			Category: "User actions",
			Name:     "login",
			Aliases:  []string{"a"},
			Usage:    "Initial sign-up and subsequent sign-in",
			Action: func(c *cli.Context) error {

				login()
				return nil
			},
		},
		{
			Category: "User actions",
			Name:     "dashboard",
			Aliases:  []string{"d"},
			Usage:    "Open dashboard at vChain.us in browser",
			Action: func(c *cli.Context) error {

				dashboard()
				return nil
			},
		},
	}

	createDirectoryInfrastructure()
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
