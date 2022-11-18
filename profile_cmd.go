package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/pterm/pterm"
	"github.com/urfave/cli/v2"
	"go.etcd.io/bbolt"
)

var profileName string = "profile"
var profileAlias []string = []string{"profiles", "pf"}

func GetProfile(DB *bbolt.DB) *cli.Command {
	return &cli.Command{
		Name:    profileName,
		Aliases: profileAlias,
		Usage:   "list of profiles",
		Action: func(ctx *cli.Context) error {

			p := &Profile{}
			ps := p.ListProfile(DB)

			var td pterm.TableData
			td = [][]string{{"NAME", "AUTH"}}
			for _, v := range ps {
				td = append(td, []string{v.Name, strconv.FormatBool(v.AskPassword)})
			}
			pterm.DefaultTable.WithHasHeader().WithData(td).Render()
			return nil
		},
	}
}

func CreateProfile(DB *bbolt.DB) *cli.Command {
	return &cli.Command{
		Name:      profileName,
		Aliases:   profileAlias,
		Usage:     "create a profile",
		UsageText: "ovpncli create profile [arguments...] [name]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "file",
				Required: true,
				Aliases:  []string{"f"},
				Usage:    "path to ovpn file",
			},
		},
		Action: func(ctx *cli.Context) error {

			// get profile name by args
			name := ctx.Args().Get(0)
			if name == "" {
				fmt.Printf("USAGE:\n   %s\n", ctx.Command.UsageText)
				return nil
			}

			// check file extension is .ovpn
			if filepath.Ext(ctx.String("file")) != ".ovpn" {
				fmt.Fprintf(os.Stderr, "\033[31munrecognize file extension '%s' supported extension ['.ovpn']\033[31m\n", filepath.Ext(ctx.String("file")))
				return nil
			}

			// read the file
			profile, err := os.ReadFile(ctx.String("file"))
			if err != nil {
				fmt.Fprintf(os.Stderr, "\033[31mfailed to parse file %s\033[31m\n", ctx.String("file"))
				return err
			}

			// check if there auth-user-pass exists or not
			var askPassword bool
			re, _ := regexp.Compile(`(?m)^auth-user-pass$`)
			isMatch := re.Match(profile)
			if isMatch {
				askPassword = true
			}

			p := &Profile{
				Name:        name,
				Profile:     profile,
				AskPassword: askPassword,
			}

			// save profile to database
			if err := p.CreateProfile(DB); err != nil {
				fmt.Fprintf(os.Stderr, "\033[31mfailed to create profile %s\033[31m\n", name)
				return err
			}

			fmt.Printf("Profile Created\n")
			return nil
		},
	}
}

func DeleteProfile(DB *bbolt.DB) *cli.Command {
	return &cli.Command{
		Name:      profileName,
		Aliases:   profileAlias,
		UsageText: "ovpncli delete profile [name]",
		Usage:     "delete a profile",
		Action: func(ctx *cli.Context) error {
			// get profile name by args
			name := ctx.Args().Get(0)
			if name == "" {
				fmt.Printf("USAGE:\n   %s\n", ctx.Command.UsageText)
				return nil
			}

			p := &Profile{
				Name: name,
			}

			// delete profile from database
			if err := p.DeleteProfile(DB); err != nil {
				fmt.Fprintf(os.Stderr, "\033[31mfailed to delete profile %s\033[31m\n", name)
				return err
			}
			fmt.Printf("Profile Deleted\n")
			return nil
		},
	}
}

func DescribeProfile(DB *bbolt.DB) *cli.Command {
	return &cli.Command{
		Name:      profileName,
		Aliases:   profileAlias,
		UsageText: "ovpncli describe profile [name]",
		Usage:     "describe a profile",
		Action: func(ctx *cli.Context) error {
			// get profile name by args
			name := ctx.Args().Get(0)
			if name == "" {
				fmt.Printf("USAGE:\n   %s\n", ctx.Command.UsageText)
				return nil
			}

			p := &Profile{
				Name: name,
			}

			res := p.DescribeProfile(DB)
			fmt.Println(string(res.Profile))
			return nil
		},
	}
}
