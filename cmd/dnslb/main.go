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
			&cli.StringFlag{
				Name:     "subdomain",
				Usage:    "Subdomain to add IPs to.",
				Required: true,
			},
			&cli.StringSliceFlag{
				Name:     "ports",
				Usage:    "Ports to check for health.",
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
						Usage: "IPv4 to use for load balancing. Options are IP, 'auto' and 'off'.",
						Value: true,
					},
					&cli.BoolFlag{
						Name:  "ipv6",
						Usage: "IPv6 to use for load balancing. Options are a IPv6, 'auto' and 'off'.",
						Value: true,
					},
				},
				Action: func(c *cli.Context) error {
					return dnslb.Endpoint(c.String("subdomain"), c.StringSlice("ports"), c.Bool("ipv4"), c.Bool("ipv6"))
				},
			},
			{
				Name:  "cleanup",
				Usage: "check endpoints and remove unhealthy entries",
				Action: func(c *cli.Context) error {
					return dnslb.Cleanup(c.String("subdomain"), c.StringSlice("ports"))
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
