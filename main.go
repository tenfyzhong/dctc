package main

import (
	"fmt"
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/urfave/cli"
)

var (
	host      string
	tls       bool
	tlscacert string
	tlscert   string
	tlskey    string
	tlsverify bool
)

func main() {
	home, err := homedir.Dir()
	capath := ""
	certpath := ""
	keypath := ""
	if err == nil {
		capath = filepath.Join(home, ".docker", "ca.pem")
		certpath = filepath.Join(home, ".docker", "cert.pem")
		keypath = filepath.Join(home, ".docker", "key.pem")
	}

	app := cli.NewApp()
	app.HideHelp = true
	app.Name = "dctc"
	app.Usage = "generate a docker-compose.yml from a container"
	app.Version = "0.1.0"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "host, H",
			Value:       "unix:///var/run/docker.sock",
			Usage:       "Daemon `socket` to connect to",
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
			Value:       capath,
			Destination: &tlscacert,
		},
		cli.StringFlag{
			Name:        "tlscert",
			Usage:       "Path to TLS certificate file",
			Value:       certpath,
			Destination: &tlscert,
		},
		cli.StringFlag{
			Name:        "tlskey",
			Usage:       "Path to TLS key file",
			Value:       keypath,
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

		str, err := dctc(client, []string(c.Args()))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return err
		}

		fmt.Println(str)

		return nil
	}

	app.Run(os.Args)
}
