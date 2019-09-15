package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/adrg/xdg"
	"github.com/esiqveland/notify"
	"github.com/godbus/dbus"
	"golang.org/x/xerrors"
)

var defaultNotificationClient = &notificationClient{
	store: &notificationStoreClient{},
}

type notificationClient struct {
	store notificationStore
}

func (c *notificationClient) showVolumeNotification(volume int, mute bool) error {
	conn, err := dbus.SessionBus()
	if err != nil {
		return xerrors.Errorf("error connecting to DBus: %w", err)
	}

	n := notify.Notification{
		AppName:       "volumectl",
		AppIcon:       c.notificationVolumeIcon(volume, mute),
		ReplacesID:    c.store.LastNoficationID(),
		ExpireTimeout: int32(3000),
		Hints: map[string]dbus.Variant{
			"value":       dbus.MakeVariant(volume),
			"synchronous": dbus.MakeVariant("volume"),
		},
	}

	notificationID, err := notify.SendNotification(conn, n)
	if err != nil {
		return xerrors.Errorf("error sending notification: %w", err)
	}
	return c.store.WriteNotificationID(notificationID)
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

type notificationStore interface {
	LastNoficationID() uint32
	WriteNotificationID(uint32) error
}

type notificationStoreClient struct{}

var _ notificationStore = &notificationStoreClient{}

func (s *notificationStoreClient) LastNoficationID() uint32 {
	storeFilePath, err := xdg.SearchRuntimeFile(s.storeFilePath())
	if err != nil {
		return 0
	}

	storeFileContent, err := ioutil.ReadFile(storeFilePath)
	if err != nil {
		return 0
	}

	lastNotificationID := strings.TrimSpace(string(storeFileContent))
	notificationID, _ := strconv.ParseUint(lastNotificationID, 10, 32)
	return uint32(notificationID)
}

func (s *notificationStoreClient) WriteNotificationID(id uint32) error {
	if err := ensureDirExists(path.Join(xdg.RuntimeDir, "volumectl")); err != nil {
		return xerrors.Errorf("unable to create XDG runtime directory: %w", err)
	}

	storeFilePath := path.Join(xdg.RuntimeDir, s.storeFilePath())
	return ioutil.WriteFile(storeFilePath, []byte(fmt.Sprintf("%d", id)), 0600)
}

func (s *notificationStoreClient) storeFilePath() string {
	return path.Join("volumectl", "last_notification_id")
}

func ensureDirExists(dir string) error {
	fi, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return os.Mkdir(dir, os.ModeDir|0700)
		}
		return err
	}

	if !fi.IsDir() {
		return xerrors.Errorf("%d is not a directory", dir)
	}
	return nil
}
