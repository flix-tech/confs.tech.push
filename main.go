package main

import (
	"log"
	"os"

	"gopkg.in/urfave/cli.v1"

	"github.com/flix-tech/confs.tech.push/cmd"
)

func main() {
	app := cli.NewApp()

	app.Name = "confs.tech.push"
	app.Usage = "push data about tech conferences to somewhere you can read it"
	app.Version = "1.2.0"

	app.Flags = []cli.Flag{
		cli.StringSliceFlag{
			Name:   "countries-blacklist, C",
			Usage:  "Countries to be blocked",
			EnvVar: "COUNTRIES_BLACKLIST",
		},
		cli.BoolFlag{
			Name:   "cfp-finished",
			Usage:  "Post only conferences with CallForPapers stage finished",
			EnvVar: "CFP_FINISHED",
		},
	}

	app.Commands = []cli.Command{
		cmd.SlackCommand(),
		cmd.AtomCommand(),
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
