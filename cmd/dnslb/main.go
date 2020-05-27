package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
	dnslb "github.com/utopia-planitia/dnslb/pkg"
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:     "port",
				Usage:    "Port to check for health. Can be defined multiple times.",
				Required: true,
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "endpoint",
				Usage: "keep healthy endpoint in dns load balancer",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "ipv4",
						Usage: "use IPv4 for load balancing.",
						Value: true,
					},
					&cli.BoolFlag{
						Name:  "ipv6",
						Usage: "use IPv6 for load balancing.",
						Value: true,
					},
				},
				Action: func(c *cli.Context) error {
					return dnslb.Endpoint(c.StringSlice("port"), c.Bool("ipv4"), c.Bool("ipv6"))
				},
			},
			{
				Name:  "cleanup",
				Usage: "check endpoints and remove unhealthy entries",
				Action: func(c *cli.Context) error {
					return dnslb.Cleanup(c.StringSlice("port"))
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
