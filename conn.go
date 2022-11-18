package main

import (
	"errors"
	"fmt"
	"os/user"
	"strings"

	"github.com/howeyc/gopass"
	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/mysteriumnetwork/go-openvpn/openvpn3"
	"github.com/urfave/cli/v2"
	"go.etcd.io/bbolt"
)

type callbacks interface {
	openvpn3.Logger
	openvpn3.EventConsumer
	openvpn3.StatsConsumer
}

type loggingCallbacks struct{}

func (lc *loggingCallbacks) Log(text string) {
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		fmt.Println("ovpncli log   >>", line)
	}
}

func (lc *loggingCallbacks) OnEvent(event openvpn3.Event) {
	fmt.Printf("ovpncli event >> %+v\n", event)
}

func (lc *loggingCallbacks) OnStats(stats openvpn3.Statistics) {
	fmt.Printf("ovpncli stats >> %+v\n", stats)
}

var _ callbacks = &loggingCallbacks{}

// StdoutLogger represents the stdout logger callback
type StdoutLogger func(text string)

// Log logs the given string to stdout logger
func (lc StdoutLogger) Log(text string) {
	lc(text)
}

func ConnectProfile(DB *bbolt.DB) *cli.Command {
	return &cli.Command{
		Name:      profileName,
		Aliases:   profileAlias,
		Usage:     "connect profile",
		UsageText: "ovpncli connect profile [name]",
		Action: func(ctx *cli.Context) error {
			// check user is root
			currentUser, err := user.Current()
			if err != nil {
				return errors.New("error check current user")
			}

			if currentUser.Username != "root" {
				fmt.Println("this command need to run by root or sudoers")
				return nil
			}

			// declare profile
			p := &Profile{}

			// get profile name by args
			name := ctx.Args().Get(0)

			// use fzf if args null
			if name == "" {
				profileList := p.ListProfile(DB)
				idx, err := fuzzyfinder.FindMulti(
					profileList,
					func(i int) string {
						return profileList[i].Name
					})
				if err != nil {
					return err
				}
				name = profileList[idx[0]].Name
			}

			p = &Profile{Name: name}
			profileResponse := p.DescribeProfile(DB)

			var logger StdoutLogger = func(text string) {
				lines := strings.Split(text, "\n")
				for _, line := range lines {
					fmt.Println("Library check >>", line)
				}
			}

			openvpn3.SelfCheck(logger)

			config := openvpn3.NewConfig(string(profileResponse.Profile))

			var username string
			var password []byte
			var userCredentials openvpn3.UserCredentials

			if profileResponse.AskPassword {
				// input username
				fmt.Print("Enter Username: ")
				fmt.Scanln(&username)

				// input password
				fmt.Print("Enter Password: ")
				password, _ = gopass.GetPasswdMasked()

				userCredentials = openvpn3.UserCredentials{Username: username, Password: string(password)}
			}

			session := openvpn3.NewSession(config, userCredentials, &loggingCallbacks{})
			session.Start()
			err = session.Wait()

			if err != nil {
				fmt.Println("Openvpn3 error: ", err)
			} else {
				fmt.Println("Graceful exit")
			}
			return nil
		},
	}
}
