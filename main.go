package main

import notify "github.com/mqu/go-notify"
import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "volumectl"
	app.Usage = "Control master volume from the command line"
	app.Version = "0.0.1"
	app.Commands = []cli.Command{
		{
			Name:  "up",
			Usage: "increase volume (with 2%)",
			Action: func(c *cli.Context) {
				exec.Command("amixer", "sset", "Master", "2%+", "unmute").Run()
				volume := getCurrentVolume()
				showVolumeNotification(volume, false)
			},
		},
		{
			Name:  "down",
			Usage: "decrease volume (with 2%)",
			Action: func(c *cli.Context) {
				exec.Command("amixer", "sset", "Master", "2%-", "unmute").Run()
				volume := getCurrentVolume()
				showVolumeNotification(volume, false)
			},
		},
		{
			Name:  "mute",
			Usage: "mute volume",
			Action: func(c *cli.Context) {
				exec.Command("amixer", "sset", "Master", "mute").Run()
				volume := getCurrentVolume()
				showVolumeNotification(volume, true)
			},
		},
		{
			Name:  "unmute",
			Usage: "unmute volume",
			Action: func(c *cli.Context) {
				exec.Command("amixer", "sset", "Master", "unmute").Run()
				volume := getCurrentVolume()
				showVolumeNotification(volume, false)
			},
		},
		{
			Name:  "toggle",
			Usage: "toggle mute",
			Action: func(c *cli.Context) {
				exec.Command("amixer", "sset", "Master", "toggle").Run()
				volume := getCurrentVolume()
				muted := mutedVolume()
				showVolumeNotification(volume, muted)
			},
		},
		{
			Name:  "set",
			Usage: "set volume to a specific value",
			Action: func(c *cli.Context) {
				exec.Command("amixer", "sset", "Master", c.Args().First()).Run()
				volume := getCurrentVolume()
				showVolumeNotification(volume, false)
			},
		},
	}
	app.Action = func(c *cli.Context) {
		volume := getCurrentVolume()
		mute := "[on]"
		if mutedVolume() {
			mute = "[off]"
		}
		fmt.Println("Volume:", strconv.Itoa(volume)+"%", mute)
	}

	app.Run(os.Args)
}

func getCurrentVolume() int {
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

	parsedVolume := string(res[1][:])
	volume, err := strconv.Atoi(parsedVolume)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to convert volume value\n")
		return 0
	}

	return volume
}

func mutedVolume() bool {
	out, err := exec.Command("amixer", "sget", "Master").Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to execute amixer command\n")
		return false
	}

	return regexp.MustCompile(`\[(\d+)%\] \[off\]`).Match(out)
}

func showVolumeNotification(volume int, mute bool) {
	notify.Init("Hello World!")
	hello := notify.NotificationNew("Hello World!",
		"This is an example notification.",
		"")

	if hello == nil {
		fmt.Fprintf(os.Stderr, "Unable to create a new notification\n")
		return
	}

	displayVolume := int32(volume)
	iconName := "notification-audio-volume-medium"
	if displayVolume == 0 {
		iconName = "notification-audio-volume-off"
	} else if displayVolume > 70 {
		iconName = "notification-audio-volume-high"
	} else if displayVolume < 30 {
		iconName = "notification-audio-volume-low"
	}

	if mute {
		iconName = "notification-audio-volume-muted"
	}

	hello.Update(" ", "", iconName)
	hello.SetHintInt32("value", displayVolume)
	hello.SetHintString("synchronous", "volume")

	if e := hello.Show(); e != nil {
		fmt.Fprintf(os.Stderr, "%s\n", e.Message())
		return
	}

	notify.UnInit()
}
