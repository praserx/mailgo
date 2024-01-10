package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/praserx/mailgo"
	"github.com/urfave/cli/v2"
)

var flags = []cli.Flag{
	&cli.StringFlag{
		Name:     "host",
		Required: true,
		Usage:    "mail server host",
	},
	&cli.IntFlag{
		Name:  "port",
		Value: 25,
		Usage: "mail server port",
	},
	&cli.StringFlag{
		Name:     "from",
		Required: true,
		Usage:    "sender address",
	},
	&cli.StringFlag{
		Name:  "name",
		Value: "MailGo Test",
		Usage: "sender name",
	},
	&cli.StringFlag{
		Name:  "username",
		Usage: "username for authentication",
	},
	&cli.StringFlag{
		Name:  "password",
		Usage: "password for authentication",
	},
}

func main() {

	app := &cli.App{
		Name:        "MailGo",
		Usage:       "Test binary for MailGo package",
		Version:     "v1.0.0",
		Description: "This application is used for testing and example purpose only.",
		Copyright:   "(c) 2023 Praser",
		Flags:       flags,
		Action: func(ctx *cli.Context) error {
			var options []mailgo.MailerOption
			options = append(options, mailgo.WithFrom(ctx.String("from")))
			options = append(options, mailgo.WithHost(ctx.String("host")))
			options = append(options, mailgo.WithPort(strconv.Itoa(ctx.Int("port"))))
			options = append(options, mailgo.WithName(ctx.String("name")))
			if ctx.String("username") != "" && ctx.String("password") != "" {
				options = append(options, mailgo.WithCredentials(ctx.String("username"), ctx.String("password")))
			}

			if mailer, err := mailgo.NewMailer(options...); err != nil {
				fmt.Fprintf(os.Stderr, "cannot initialize mailer: %v", err.Error())
				os.Exit(1)
			} else {
				errs := mailer.SendMail([]string{"praserx@gmail.com"}, "Test mail", "Test", "<p>Test</p>")
				for _, err := range errs {
					if err != nil {
						fmt.Fprintf(os.Stderr, "fatal: %v", err.Error())
					}
				}

				if len(errs) != 0 {
					os.Exit(1)
				}

			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "fatal error: %s\n", err.Error())
		os.Exit(1)
	}
}
