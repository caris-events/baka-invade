package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "baka-invade",
		Usage: "only baka invades the others.",
		Commands: []*cli.Command{
			{
				Name:  "update",
				Usage: "updates the necessary objects and dicts",
				Action: func(c *cli.Context) error {
					return Update(c)
				},
			},
			{
				Name:  "refresh",
				Usage: "refresh the template files like: index.html, about.html and styles",
				Action: func(c *cli.Context) error {
					return Refresh(c)
				},
			},
			{
				Name:  "rebuild",
				Usage: "clear the cache and rebuild all objects, dicts and cover images (it will take a while!)",
				Action: func(c *cli.Context) error {
					return Rebuild(c)
				},
			},
			// {
			// 	Name:    "export",
			// 	Aliases: []string{"c"},
			// 	Usage:   "compiles and export the data as json files",
			// 	Action: func(cCtx *cli.Context) error {
			// 		return nil
			// 	},
			// },
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
