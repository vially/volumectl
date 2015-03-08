package main

import (
	"log"

	"github.com/mqu/go-notify"
)

func showVolumeNotification(volume int, mute bool) {
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
	notification.Show()

	notify.UnInit()
}
