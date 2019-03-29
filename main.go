package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

var (
	debug     bool
	file      string
	host      string
	tls       bool
	tlscacert string
	tlscert   string
	tlskey    string
	tlsverify bool
	version   string
)

func main() {
	app := cli.NewApp()
	app.Name = "dctc"
	app.Usage = "generate a docker-compose.yml from a container"
	app.Version = "0.1.0"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "debug",
			Usage:       "Enable debug mode",
			Destination: &debug,
		},
		cli.StringFlag{
			Name:        "host, H",
			Value:       "unix:///var/run/docker.sock",
			Usage:       "Daemon socket to connect to",
			Destination: &host,
		},
		cli.BoolFlag{
			Name:        "tls",
			Usage:       "Use TLS; implied by --tlsverify",
			Destination: &tls,
		},
		cli.StringFlag{
			Name:        "tlscacert",
			Usage:       "Trust certs signed only by this CA",
			Destination: &tlscacert,
		},
		cli.StringFlag{
			Name:        "tlscert",
			Usage:       "Path to TLS certificate file",
			Destination: &tlscert,
		},
		cli.StringFlag{
			Name:        "tlskey",
			Usage:       "Path to TLS key file",
			Destination: &tlskey,
		},
		cli.BoolFlag{
			Name:        "tlsverify",
			Usage:       "Use TLS and verify the remote",
			Destination: &tlsverify,
		},
	}

	app.Action = func(c *cli.Context) error {
		if c.NArg() == 0 {
			fmt.Fprintf(os.Stderr, "need an argument\n\n")
			cli.ShowAppHelpAndExit(c, -1)
		}

		client, err := newClient()
		if err != nil {
			return err
		}
		err = dctc(client, c.Args().Get(0))
		if err != nil {
			return err
		}
		return nil
	}

	app.Run(os.Args)
}
