// Copyright (C) 2015 Space Monkey, Inc.

package sdl

import (
	"path/filepath"
	"strings"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_mixer"

	"sm/final/assets"
)

func cleanupSounds(sounds map[string]*mix.Chunk) {
	cleaned_sounds := map[*mix.Chunk]bool{}
	for _, sound := range sounds {
		if !cleaned_sounds[sound] {
			sound.Free()
			cleaned_sounds[sound] = true
		}
	}
}

func BackgroundMusic(asset string) error {
	return exec(func() error {
		data, err := assets.Asset(asset)
		if err != nil {
			return err
		}
		mem := sdl.RWFromMem(unsafe.Pointer(&data[0]), len(data))
		music, err := mix.LoadMUS_RW(mem, 0)
		if err != nil {
			return err
		}
		return music.Play(-1)
	})
}

func loadSounds() (
	sounds map[string]*mix.Chunk, err error) {
	files, err := assets.AssetDir("final/sounds")
	if err != nil {
		return nil, err
	}
	sounds = make(map[string]*mix.Chunk)
	for _, file := range files {
		if !strings.HasSuffix(file, ".wav") {
			continue
		}
		data, err := assets.Asset(filepath.Join("final/sounds", file))
		if err != nil {
			return nil, err
		}
		mem := sdl.RWFromMem(unsafe.Pointer(&data[0]), len(data))
		sound, err := mix.LoadWAV_RW(mem, false)
		if err != nil {
			cleanupSounds(sounds)
			return nil, err
		}
		sounds[file[:len(file)-len(".wav")]] = sound
	}
	return sounds, nil
}
