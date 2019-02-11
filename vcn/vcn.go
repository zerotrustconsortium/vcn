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

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "vcn"
	app.Version = "0.0.1"

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
			Usage:    "Verify against blockchain",
			Action: func(c *cli.Context) error {
				if c.NArg() == 0 {
					return fmt.Errorf("filenames required")
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
			Name:     "commit",
			Aliases:  []string{"c"},
			Usage:    "Commit in blockchain",
			Action: func(c *cli.Context) error {
				if c.NArg() == 0 {
					return fmt.Errorf("filename required")
				}
				Commit(c.Args().First(), "me")
				return nil
			},
		},
		{
			Category: "User actions",
			Name:     "auth",
			Aliases:  []string{"a"},
			Usage:    "Authenticate against vChain.us.",
			Action: func(c *cli.Context) error {

				Authenticate()
				return nil
			},
			Flags: []cli.Flag{
				cli.BoolFlag{Name: "hash"},
			},
		},
		{
			Category: "User actions",
			Name:     "register",
			Aliases:  []string{"r"},
			Usage:    "Register an account on the vChain.us platform.",
			Action: func(c *cli.Context) error {

				Register()
				return nil
			},
			Flags: []cli.Flag{
				cli.BoolFlag{Name: "hash"},
			},
		},
	}

	createDirectoryInfrastructure()
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
