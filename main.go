package main

import (
	"fmt"
	"os"

	"github.com/phayes/freeport"
	"github.com/urfave/cli/v2"
)

func main() {
	var app = cli.NewApp()

	app.Name = "monocular"
	app.Usage = "view the logs on the server in real time through the web page"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "listen-address",
			Usage:       "service listen address",
			DefaultText: "127.0.0.1 and random port",
		},
	}
	app.ArgsUsage = "some.log"
	app.HideHelpCommand = true
	app.Action = func(ctx *cli.Context) error {
		var listenAddress = ctx.String("listen-address")
		if listenAddress == "" {
			port, err := freeport.GetFreePort()
			if err != nil {
				return err
			}
			listenAddress = fmt.Sprintf("127.0.0.1:%d", port)
		}

		var filepath = ctx.Args().First()
		if filepath == "" {
			return fmt.Errorf("please specify log file")
		}

		return startServer(listenAddress, filepath)
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err.Error())
	}
}
