package main

import (
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"strconv"
	"syncbear/ray"
	"syncbear/ss"
	"time"
)

var users = make(map[string]string)

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "mode",
			Value: "v2ray",
			Usage: "backend mode:shadowsocks v2ray trojan",
		},
		&cli.StringFlag{
			Name:     "url",
			Value:    "www.yourdomain.com",
			Usage:    "server address for sync user info and traffic, no https://",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "token",
			Value:    "yourtoken",
			Usage:    "your token to access api service",
			Required: true,
		},
		&cli.IntFlag{
			Name:  "port",
			Value: 8080,
			Usage: "manage api port",
		},
		&cli.IntFlag{
			Name:  "sync",
			Value: 60,
			Usage: "sync user info interval",
		},
		&cli.IntFlag{
			Name:  "traffic",
			Value: 600,
			Usage: "traffic update interval",
		},
	}

	app.Action = func(c *cli.Context) error {
		mode := c.String("mode")
		port := c.Int("port")
		pstr := strconv.Itoa(port)
		tick1 := time.Tick(time.Duration(c.Int("sync")) * time.Second)
		tick2 := time.Tick(time.Duration(c.Int("traffic")) * time.Second)

		if mode == "shadowsocks" {
			client, err := ss.New("127.0.0.1:" + pstr)
			if err != nil {
				return err
			}

			for {
				select {
				case <-tick1:
					go ss.Sync(client, users, c.String("url"), c.String("token"))
				case <-tick2:
					go ss.SyncTraffic(client, c.String("url"), c.String("token"))
				}
			}

		}

		if mode == "v2ray" {
			client := ray.New("127.0.0.1", c.Int("port"))
			for {
				select {
				case <-tick1:
					go ray.Sync(client, users, c.String("url"), c.String("token"))
				case <-tick2:
					go ray.Traffic(client, c.String("url"), c.String("token"))
				}
			}
		}
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
