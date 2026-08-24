package utils

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

var audioContext = audio.NewContext(24000)

// Play a sound.
func PlaySound(name string) {
	// Load Stream.
	data := LoadWav(name)

	// Create a player
	player, err := audioContext.NewPlayerF32(data)
	if err != nil {
		log.Fatal(err)
	}

	// Set Volume.
	switch name {
	case "jump":
		player.SetVolume(0.2)
	case "stomp":
		player.SetVolume(0.2)
	case "kick":
		player.SetVolume(0.3)
	case "died":
		player.SetVolume(0.3)
	case "shrink":
		player.SetVolume(0.4)
	case "coin":
		player.SetVolume(0.1)
	case "brick_break":
		player.SetVolume(0.3)
	case "1_up":
		player.SetVolume(0.3)
	case "game_over":
		player.SetVolume(0.6)
	case "stage_clear":
		player.SetVolume(0.3)
	}

	// Play the audio.
	player.Play()
}

// Play an audio on loop.
func PlayLoop(name string) (player *audio.Player) {
	// Load Stream.
	data := LoadWav(name)

	// Create a loop.
	loop := audio.NewInfiniteLoop(data, data.Length())

	// Create a player.
	player, err := audioContext.NewPlayerF32(loop)
	if err != nil {
		log.Fatal(err)
	}

	// Set Volume.
	switch name {
	case "music":
		player.SetVolume(0.1)
	}

	// Play the audio.
	player.Play()

	return player
}

// Load a Wav stream from the given name.
func LoadWav(name string) *wav.Stream {
	// Open file.
	file, err := Assets.Sfx.Open(fmt.Sprintf("sfx/%s.wav", name))
	if err != nil {
		log.Fatal(err)
	}

	// Decode wav.
	data, err := wav.DecodeF32(file)
	if err != nil {
		log.Fatal(err)
	}

	return data
}
