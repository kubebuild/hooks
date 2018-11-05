package main

import (
	"os"

	webhook "github.com/kubebuild/webhooks/pkg/app"
	"github.com/urfave/cli"
)

var appHelpTemplate = `{{.Name}} - {{.Usage}}

OPTIONS:
  {{range .Flags}}{{.}}
  {{end}}
`

func main() {
	cli.AppHelpTemplate = appHelpTemplate

	app := cli.NewApp()

	app.Name = "webhoos"
	app.Version = "1.0.0"
	app.Usage = "kubebuild webhooks and validator"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "graphql-url",
			Value: "https://api.kubebuild.com/graphql",
			Usage: "api url for graphql",
		},
		cli.StringFlag{
			Name:  "log-level",
			Value: "info",
			Usage: "log level",
		},
	}
	app.Action = func(c *cli.Context) {
		logLevel := c.String("log-level")
		version := c.App.Version
		name := c.App.Name
		graphqlURL := c.String("graphql-url")
		config := &webhook.Config{
			Name:       name,
			Version:    version,
			GraphqlURL: graphqlURL,
			LogLevel:   logLevel,
		}
		_, err := webhook.NewApp(config)
		if err != nil {
			panic("Error occured exiting...")
		}
	}
	app.Run(os.Args)
}
