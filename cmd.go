package main

import (
	"log"
	"os"
	"runtime"

	"github.com/urfave/cli/v2"
)

func main() {
	// 默认并发数
	concurrencyN := runtime.NumCPU()

	app := &cli.App{
		Name:  "fastfetch",
		Usage: "File concurrency download",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "url",
				Aliases:  []string{"u"},
				Usage:    "`URL` to download",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "Output `filename`",
			},
			&cli.IntFlag{
				Name:    "max-connect",
				Aliases: []string{"n"},
				Value:   concurrencyN,
				Usage:   "Specify maximum number of connections",
			},
		},
		Action: func(c *cli.Context) error {
			strURL := c.String("url")
			filename := c.String("output")
			concurrency := c.Int("max-connect")
			return NewDownloader(concurrency).Download(strURL, filename, concurrency)
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
