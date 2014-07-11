package main

import notify "github.com/mqu/go-notify"
import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
	"os/exec"
	"regexp"
	"strconv"
)

func main() {
	app := cli.NewApp()
	app.Name = "volumectl"
	app.Usage = "Control master volume from the command line"
	app.Version = "0.0.1"
	app.Commands = []cli.Command{
		{
			Name:  "up",
			Usage: "increase volume",
			Action: func(c *cli.Context) {
				exec.Command("amixer", "sset", "Master", "2%+", "unmute").Run()
				volume := GetCurrentVolume()
				ShowVolumeNotification(volume, false)
			},
		},
		{
			Name:  "down",
			Usage: "decrease volume",
			Action: func(c *cli.Context) {
				exec.Command("amixer", "sset", "Master", "2%-", "unmute").Run()
				volume := GetCurrentVolume()
				ShowVolumeNotification(volume, false)
			},
		},
		{
			Name:  "mute",
			Usage: "mute volume",
			Action: func(c *cli.Context) {
				exec.Command("amixer", "sset", "Master", "mute").Run()
				volume := GetCurrentVolume()
				ShowVolumeNotification(volume, true)
			},
		},
		{
			Name:  "unmute",
			Usage: "unmute volume",
			Action: func(c *cli.Context) {
				exec.Command("amixer", "sset", "Master", "unmute").Run()
				volume := GetCurrentVolume()
				ShowVolumeNotification(volume, false)
			},
		},
		{
			Name:  "toggle",
			Usage: "toggle mute",
			Action: func(c *cli.Context) {
				exec.Command("amixer", "sset", "Master", "toggle").Run()
				volume := GetCurrentVolume()
				muted := MutedVolume()
				ShowVolumeNotification(volume, muted)
			},
		},
		{
			Name:  "set",
			Usage: "set volume",
			Action: func(c *cli.Context) {
				exec.Command("amixer", "sset", "Master", c.Args().First()).Run()
				volume := GetCurrentVolume()
				ShowVolumeNotification(volume, false)
			},
		},
	}
	app.Action = func(c *cli.Context) {
		volume := GetCurrentVolume()
		mute := "[on]"
		if MutedVolume() {
			mute = "[off]"
		}
		fmt.Println("Volume:", strconv.Itoa(volume)+"%", mute)
	}

	app.Run(os.Args)
}

func GetCurrentVolume() int {
	out, err := exec.Command("amixer", "sget", "Master").Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to execute amixer command\n")
		return 0
	}

	re := regexp.MustCompile(`\[(\d+)%\]`)
	res := re.FindSubmatch(out)
	if res == nil {
		fmt.Fprintf(os.Stderr, "Unable to parse volume from amixer command output\n")
		return 0
	}

	parsed_volume := string(res[1][:])
	volume, err := strconv.Atoi(parsed_volume)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to convert volume value\n")
		return 0
	}

	return volume
}

func MutedVolume() bool {
	out, err := exec.Command("amixer", "sget", "Master").Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to execute amixer command\n")
		return false
	}

	return regexp.MustCompile(`\[(\d+)%\] \[off\]`).Match(out)
}

func ShowVolumeNotification(volume int, mute bool) {
	notify.Init("Hello World!")
	hello := notify.NotificationNew("Hello World!",
		"This is an example notification.",
		"")

	if hello == nil {
		fmt.Fprintf(os.Stderr, "Unable to create a new notification\n")
		return
	}

	display_volume := int32(volume)
	icon_name := "notification-audio-volume-medium"
	if display_volume == 0 {
		icon_name = "notification-audio-volume-off"
	} else if display_volume > 70 {
		icon_name = "notification-audio-volume-high"
	} else if display_volume < 30 {
		icon_name = "notification-audio-volume-low"
	}

	if mute {
		icon_name = "notification-audio-volume-muted"
	}

	hello.Update(" ", "", icon_name)
	hello.SetHintInt32("value", display_volume)
	hello.SetHintString("synchronous", "volume")

	if e := hello.Show(); e != nil {
		fmt.Fprintf(os.Stderr, "%s\n", e.Message())
		return
	}

	notify.UnInit()
}
