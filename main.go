package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "rollbar-drone"
	app.Usage = "rollbar drone plugin"
	app.Action = run
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "access-token",
			Usage:  "rolbar access token",
			EnvVar: "PLUGIN_ROLLBAR_ACCESS_TOKEN",
		},
		cli.StringFlag{
			Name:   "environment",
			Usage:  "rolbar environment",
			EnvVar: "PLUGIN_ROLLBAR_ENVIRONMENT",
		},
		cli.StringFlag{
			Name:   "revision",
			Usage:  "git commit sha",
			EnvVar: "DRONE_COMMIT_SHA",
			Value:  "00000000",
		},
		cli.StringFlag{
			Name:   "author",
			Usage:  "git author username",
			EnvVar: "DRONE_COMMIT_AUTHOR",
		},
		cli.StringFlag{
			Name:   "comment",
			Usage:  "commit message",
			EnvVar: "DRONE_COMMIT_MESSAGE",
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	plugin := Plugin{
		AccessToken:   c.String("access-token"),
		Environment:   c.String("environment"),
		Revision:      c.String("revision"),
		LocalUsername: c.String("author"),
		Comment:       c.String("comment"),
	}

	return plugin.Exec()
}

type Plugin struct {
	AccessToken   string `json:"access_token"`
	Environment   string `json:"environment"`
	Revision      string `json:"revision"`
	LocalUsername string `json:"local_username"`
	Comment       string `json:"comment"`
}

func (p Plugin) Exec() error {
	body, err := json.Marshal(p)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://api.rollbar.com/api/1/deploy/", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
