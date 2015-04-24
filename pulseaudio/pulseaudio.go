package pulseaudio

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// PulseAudio client
type PulseAudio struct {
	defaultSink string
	Volume      int
	Muted       bool
}

// New returns a new PulseAudio client
func New() *PulseAudio {
	pa := &PulseAudio{defaultSink: detectDefaultSink()}
	pa.detectCurrentVolume()
	return pa
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
	pa.detectCurrentVolume()
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

func detectDefaultSink() string {
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

func (pa *PulseAudio) detectCurrentVolume() {
	out, err := exec.Command("pactl", "list", "sinks").Output()
	if err != nil {
		log.Fatalf("pactl get volume command failed: %s", err)
	}

	currentSinkOutput, err := findSinkByName(out, pa.defaultSink)
	if err != nil {
		log.Fatal(err)
	}

	re := regexp.MustCompile(`(\d+)%`)
	res := re.FindSubmatch(currentSinkOutput)
	if res == nil {
		log.Fatal("Unable to parse volume from pactl command output")
	}

	parsedVolume := string(res[1][:])
	pa.Volume, err = strconv.Atoi(parsedVolume)
	if err != nil {
		log.Fatal("Unable to convert volume value")
	}
	pa.Muted = regexp.MustCompile(`Mute: yes`).Match(currentSinkOutput)
}

func findSinkByName(output []byte, sinkName string) ([]byte, error) {
	sinks := strings.Split(string(output), "Sink #")
	for _, sinkOutput := range sinks {
		if strings.Contains(sinkOutput, fmt.Sprintf("Name: %s", sinkName)) {
			return []byte(sinkOutput), nil
		}
	}
	return nil, fmt.Errorf("Unable to find sink named %s in 'pactl list sinks` output", sinkName)
}
