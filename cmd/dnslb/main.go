package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
	dnslb "github.com/utopia-planitia/dnslb/pkg"
)

func main() {
	app := &cli.App{
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
					&cli.StringSliceFlag{
						Name:     "port",
						Usage:    "Port to check for health. Can be defined multiple times.",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					return dnslb.Endpoint(c.StringSlice("port"), c.Bool("ipv4"), c.Bool("ipv6"))
				},
			},
			{
				Name:  "cleanup",
				Usage: "Check DNS endpoints and remove unhealthy entries.",
				Flags: []cli.Flag{
					&cli.StringSliceFlag{
						Name:     "port",
						Usage:    "Port to check for health. Can be defined multiple times.",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					return dnslb.Cleanup(c.StringSlice("port"))
				},
			},
			{
				Name:  "ipv4",
				Usage: "Print the detected local IPv4 address.",
				Action: func(c *cli.Context) error {
					return dnslb.IPv4()
				},
			},
			{
				Name:  "ipv6",
				Usage: "Print the detected local IPv6 address.",
				Action: func(c *cli.Context) error {
					return dnslb.IPv6()
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
