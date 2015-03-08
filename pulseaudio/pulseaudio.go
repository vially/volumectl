package pulseaudio

import (
	"log"
	"os/exec"
	"regexp"
	"strconv"
)

// PulseAudio client
type PulseAudio struct {
	defaultSink string
	Volume      int
	Muted       bool
}

// New returns a new PulseAudio client
func New() *PulseAudio {
	volume, muted := getCurrentVolume()
	return &PulseAudio{defaultSink: getDefaultSink(), Volume: volume, Muted: muted}
}

// SetMute mutes or unmutes the default sink
func (pa *PulseAudio) SetMute(muted bool) {
	mute := "0"
	if muted {
		mute = "1"
	}
	pa.Muted = muted
	exec.Command("pactl", "set-sink-mute", pa.defaultSink, mute).Run()
}

// ToggleMute toogles mute on the default sink
func (pa *PulseAudio) ToggleMute() {
	pa.Muted = !pa.Muted
	exec.Command("pactl", "set-sink-mute", pa.defaultSink, "toggle").Run()
}

// SetVolume sets the volume on the default sink
func (pa *PulseAudio) SetVolume(volume string) {
	exec.Command("pactl", "set-sink-volume", pa.defaultSink, volume).Run()
}

// IncreaseVolume increases the volume on the default sink by 2%
func (pa *PulseAudio) IncreaseVolume() {
	pa.SetMute(false)
	volumeValue := "+2%"
	pa.Volume += 2
	if pa.Volume >= 98 {
		volumeValue = "100%"
		pa.Volume = 100
	}
	exec.Command("pactl", "set-sink-volume", pa.defaultSink, volumeValue).Run()
}

// DecreaseVolume decreases the volume on the default sink by 2%
func (pa *PulseAudio) DecreaseVolume() {
	pa.SetMute(false)
	pa.Volume -= 2
	if pa.Volume < 0 {
		pa.Volume = 0
	}
	exec.Command("pactl", "set-sink-volume", pa.defaultSink, "-2%").Run()
}

func getDefaultSink() string {
	out, err := exec.Command("pactl", "info").Output()
	if err != nil {
		log.Fatalf("Unable to detect default sink: %s", err)
	}

	re := regexp.MustCompile(`Default Sink: (.*)`)
	res := re.FindSubmatch(out)
	if res == nil {
		log.Fatal("Unable to parse default sink name from pactl command output")
	}

	defaultSink := string(res[1][:])
	return defaultSink
}

func getCurrentVolume() (int, bool) {
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

	return volume, regexp.MustCompile(`Mute: yes`).Match(out)
}
