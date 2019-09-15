package main

import (
	"github.com/esiqveland/notify"
	"github.com/godbus/dbus"
	"golang.org/x/xerrors"
)

var defaultNotificationClient = &notificationClient{}

type notificationClient struct{}

func (c *notificationClient) showVolumeNotification(volume int, mute bool) error {
	conn, err := dbus.SessionBus()
	if err != nil {
		return xerrors.Errorf("error connecting to DBus: %w", err)
	}

	n := notify.Notification{
		AppName:       "volumectl",
		AppIcon:       c.notificationVolumeIcon(volume, mute),
		ExpireTimeout: int32(3000),
		Hints: map[string]dbus.Variant{
			"value":       dbus.MakeVariant(volume),
			"synchronous": dbus.MakeVariant("volume"),
		},
	}

	if _, err := notify.SendNotification(conn, n); err != nil {
		return xerrors.Errorf("error sending notification: %w", err)
	}
	return nil
}

func (c *notificationClient) notificationVolumeIcon(volume int, mute bool) string {
	iconName := "notification-audio-volume-medium"
	if mute {
		iconName = "notification-audio-volume-muted"
	} else if volume == 0 {
		iconName = "notification-audio-volume-off"
	} else if volume > 70 {
		iconName = "notification-audio-volume-high"
	} else if volume < 30 {
		iconName = "notification-audio-volume-low"
	}
	return iconName
}

func showVolumeNotification(volume int, mute bool) error {
	return defaultNotificationClient.showVolumeNotification(volume, mute)
}
