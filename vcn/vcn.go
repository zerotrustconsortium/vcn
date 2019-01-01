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
		{
			Category: "Artifact actions",
			Name:     "verify",
			Aliases:  []string{"v"},
			Usage:    "verify against blockchain",
			Action: func(c *cli.Context) error {
				fmt.Println("verified artifact: ", c.Args().First())
				verify(c.Args().First())
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
			Usage:    "commit in blockchain",
			Action: func(c *cli.Context) error {
				fmt.Println("committed artifact: ", c.Args().First())
				commit(c.Args().First(), c.Args().Get(1))
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
