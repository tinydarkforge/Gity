package services

import (
	"os/exec"
	"runtime"
)

type SoundEvent int

const (
	SoundTick  SoundEvent = iota // AI finished a turn
	SoundDone                    // issue created
	SoundError                   // something went wrong
)

var macSounds = map[SoundEvent]string{
	SoundTick:  "/System/Library/Sounds/Tink.aiff",
	SoundDone:  "/System/Library/Sounds/Hero.aiff",
	SoundError: "/System/Library/Sounds/Basso.aiff",
}

// Play emits an async terminal sound. It is a no-op when:
//   - enabled is false
//   - the platform has no supported sound command
//   - the sound binary is not found
func Play(enabled bool, event SoundEvent) {
	if !enabled {
		return
	}
	go func() {
		switch runtime.GOOS {
		case "darwin":
			file, ok := macSounds[event]
			if !ok {
				return
			}
			_ = exec.Command("afplay", file).Run()
		case "linux":
			for _, cmd := range []string{"paplay", "aplay"} {
				if path, err := exec.LookPath(cmd); err == nil {
					_ = exec.Command(path, "/usr/share/sounds/freedesktop/stereo/bell.oga").Run()
					return
				}
			}
		}
	}()
}
