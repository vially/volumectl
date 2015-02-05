package main

import (
	"log"

	notify "github.com/mqu/go-notify"
)
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
	app.Version = "0.0.3"
	app.Commands = []cli.Command{
		{
			Name:  "up",
			Usage: "increase volume (with 2%)",
			Action: func(c *cli.Context) {
				exec.Command("pactl", "set-sink-mute", "0", "0").Run()
				exec.Command("pactl", "set-sink-volume", "0", "+2%").Run()
				showVolumeNotification(false)
			},
		},
		{
			Name:  "down",
			Usage: "decrease volume (with 2%)",
			Action: func(c *cli.Context) {
				exec.Command("pactl", "set-sink-mute", "0", "0").Run()
				exec.Command("pactl", "set-sink-volume", "0", "-2%").Run()
				showVolumeNotification(false)
			},
		},
		{
			Name:  "mute",
			Usage: "mute volume",
			Action: func(c *cli.Context) {
				exec.Command("pactl", "set-sink-mute", "0", "1").Run()
				showVolumeNotification(true)
			},
		},
		{
			Name:  "unmute",
			Usage: "unmute volume",
			Action: func(c *cli.Context) {
				exec.Command("pactl", "set-sink-mute", "0", "0").Run()
				showVolumeNotification(false)
			},
		},
		{
			Name:  "toggle",
			Usage: "toggle mute",
			Action: func(c *cli.Context) {
				exec.Command("pactl", "set-sink-mute", "0", "toggle").Run()
				muted := mutedVolume()
				showVolumeNotification(muted)
			},
		},
		{
			Name:  "set",
			Usage: "set volume to a specific value",
			Action: func(c *cli.Context) {
				exec.Command("pactl", "set-sink-volume", "0", c.Args().First()).Run()
				muted := mutedVolume()
				showVolumeNotification(muted)
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
	out, err := exec.Command("pactl", "list", "sinks").Output()
	if err != nil {
		log.Fatalf("pactl get volume command failed: %s", err)
	}

	re := regexp.MustCompile(`(\d+)%`)
	res := re.FindSubmatch(out)
	if res == nil {
		log.Fatal("Unable to parse volume from pactl command output")
	}

	parsedVolume := string(res[1][:])
	volume, err := strconv.Atoi(parsedVolume)
	if err != nil {
		log.Fatal("Unable to convert volume value")
	}

	return volume
}

func mutedVolume() bool {
	out, err := exec.Command("pactl", "list", "sinks").Output()
	if err != nil {
		log.Fatalf("pactl list sinks command failed: %s", err)
	}

	return regexp.MustCompile(`Mute: yes`).Match(out)
}

func showVolumeNotification(mute bool) {
	volume := getCurrentVolume()

	notify.Init("volumectl")
	notification := notify.NotificationNew("volumectl", "This is an example notification.", "")

	if notification == nil {
		log.Fatal("Unable to create a new notification")
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

	notification.Update(" ", "", iconName)
	notification.SetHintInt32("value", displayVolume)
	notification.SetHintString("synchronous", "volume")

	if err := notification.Show(); err != nil {
		log.Fatalf("Unable to display notification. %s", err.Message())
	}

	notify.UnInit()
}
