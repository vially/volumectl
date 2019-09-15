package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/urfave/cli"
	"github.com/vially/volumectl/pulseaudio"
)

func main() {
	pa := pulseaudio.New()

	app := cli.NewApp()
	app.Name = "volumectl"
	app.Usage = "Control volume from the command line"
	app.Version = "0.0.5"
	app.Commands = []cli.Command{
		{
			Name:  "up",
			Usage: "increase volume (with 2%)",
			Action: func(c *cli.Context) {
				pa.IncreaseVolume()
				showVolumeNotification(pa.Volume, pa.Muted)
			},
		},
		{
			Name:  "down",
			Usage: "decrease volume (with 2%)",
			Action: func(c *cli.Context) {
				pa.DecreaseVolume()
				showVolumeNotification(pa.Volume, pa.Muted)
			},
		},
		{
			Name:  "mute",
			Usage: "mute volume",
			Action: func(c *cli.Context) {
				pa.SetMute(true)
				showVolumeNotification(pa.Volume, pa.Muted)
			},
		},
		{
			Name:  "unmute",
			Usage: "unmute volume",
			Action: func(c *cli.Context) {
				pa.SetMute(false)
				showVolumeNotification(pa.Volume, pa.Muted)
			},
		},
		{
			Name:  "toggle",
			Usage: "toggle mute",
			Action: func(c *cli.Context) {
				pa.ToggleMute()
				showVolumeNotification(pa.Volume, pa.Muted)
			},
		},
		{
			Name:  "set",
			Usage: "set volume to a specific value",
			Action: func(c *cli.Context) {
				pa.SetVolume(c.Args().First())
				showVolumeNotification(pa.Volume, pa.Muted)
			},
		},
	}
	app.Action = func(c *cli.Context) {
		mute := "[on]"
		if pa.Muted {
			mute = "[off]"
		}
		fmt.Println("Volume:", strconv.Itoa(pa.Volume)+"%", mute)
	}

	app.Run(os.Args)
}
