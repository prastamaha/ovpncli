package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {

	// init database
	dbpath := os.Getenv("OVPNCLI_DATABASE_PATH")
	if dbpath == "" {
		fmt.Println("Plase set environment variable 'OVPNCLI_DATABASE_PATH' within /etc/environment")
		return
	}

	DB := InitDatabase(dbpath)
	defer DB.Close()

	// init cli app
	app := &cli.App{
		Name:                   "ovpncli",
		Usage:                  "manage openvpn profiles",
		Version:                "v0.1.0",
		EnableBashCompletion:   true,
		Suggest:                true,
		UseShortOptionHandling: true,
		Commands: []*cli.Command{
			{
				Name:  "get",
				Usage: "get resource",
				Subcommands: []*cli.Command{
					GetProfile(DB),
				},
			},
			{
				Name:  "create",
				Usage: "create resource",
				Subcommands: []*cli.Command{
					CreateProfile(DB),
				},
			},
			{
				Name:  "delete",
				Usage: "delete resource",
				Subcommands: []*cli.Command{
					DeleteProfile(DB),
				},
			},
			{
				Name:  "describe",
				Usage: "describe resource",
				Subcommands: []*cli.Command{
					DescribeProfile(DB),
				},
			},
			{
				Name:  "connect",
				Usage: "connect resource",
				Subcommands: []*cli.Command{
					ConnectProfile(DB),
				},
			},
		},
	}

	// run cli app
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}
