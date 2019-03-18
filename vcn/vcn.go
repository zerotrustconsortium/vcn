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
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

// VCN_VERSION sets the version for the build + some logging with analytics
var VCN_VERSION = "0.3.3"

// WG waitgroup for sync of threads across the whole project
var WG sync.WaitGroup

func main() {
	InitLogging()

	var publicSigning = false

	LOG.WithFields(logrus.Fields{
		"version": VCN_VERSION,
	}).Trace("Started vcn")

	app := cli.NewApp()
	app.Name = "vcn"
	app.Usage = "code signing made easy"
	app.Version = VCN_VERSION

	app.Commands = []cli.Command{

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
			Usage:    "Sign digital assets' hashes onto the blockchain",
			Action: func(c *cli.Context) error {
				if c.NArg() == 0 {
					return fmt.Errorf("filename or type:reference required")
				}
				Sign(c.Args().First(), TRUSTED, visibilityForFlag(publicSigning))
				return nil
			},
			Flags: []cli.Flag{
				cli.BoolFlag{Name: "public", Destination: &publicSigning},
			},
		},
		{
			Category: "Artifact actions",
			Name:     "untrust",
			Aliases:  []string{"ut"},
			Usage:    "Untrust a digital asset.",
			Action: func(c *cli.Context) error {
				if c.NArg() == 0 {
					return fmt.Errorf("filename or type:reference required")
				}
				Sign(c.Args().First(), UNTRUSTED, visibilityForFlag(publicSigning))
				return nil
			},
			Flags: []cli.Flag{
				cli.BoolFlag{Name: "public", Destination: &publicSigning},
			},
		},
		{
			Category: "Artifact actions",
			Name:     "unsupport",
			Aliases:  []string{"ut"},
			Usage:    "Unsupport a digital asset.",
			Action: func(c *cli.Context) error {
				if c.NArg() == 0 {
					return fmt.Errorf("filename or type:reference required")
				}
				Sign(c.Args().First(), UNSUPPORTED, visibilityForFlag(publicSigning))
				return nil
			},
			Flags: []cli.Flag{
				cli.BoolFlag{Name: "public", Destination: &publicSigning},
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
			Usage:    "Sign-in to vChain.us",
			Action: func(c *cli.Context) error {

				login(nil)
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
