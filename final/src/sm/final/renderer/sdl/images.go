// Copyright (C) 2015 Space Monkey, Inc.

package sdl

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"unsafe"

	"code.google.com/p/graphics-go/graphics"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_image"

	"sm/final/assets"
	"sm/final/grid"
)

const (
	tau = 6.28318530718
)

func cleanupImages(surfaces map[grid.Cell]*sdl.Surface) {
	cleaned := map[*sdl.Surface]bool{}
	for _, surf := range surfaces {
		if !cleaned[surf] {
			surf.Free()
			cleaned[surf] = true
		}
	}
}

func cleanupRotations(surfaces map[grid.Owner][]*sdl.Surface) {
	cleaned := map[*sdl.Surface]bool{}
	for _, surf_list := range surfaces {
		for _, surf := range surf_list {
			if !cleaned[surf] {
				surf.Free()
				cleaned[surf] = true
			}
		}
	}
}

func loadImages(frames, players int) (
	images map[grid.Cell]*sdl.Surface,
	rotations map[grid.Owner][]*sdl.Surface, err error) {
	images = make(map[grid.Cell]*sdl.Surface)
	for _, simple := range []string{"battery", "floor", "wall"} {
		for _, exploding := range []bool{false, true} {
			data, err := assets.Asset(fmt.Sprintf("final/images/%s.png", simple))
			if err != nil {
				cleanupImages(images)
				return nil, nil, err
			}
			mem := sdl.RWFromMem(unsafe.Pointer(&data[0]), len(data))
			surface, err := img.Load_RW(mem, 0)
			if err != nil {
				cleanupImages(images)
				return nil, nil, err
			}
			if exploding {
				err = overlay(surface, "final/images/ex.png")
				if err != nil {
					surface.Free()
					cleanupImages(images)
					return nil, nil, err
				}
			}
			for orientation := 0; orientation < 4; orientation++ {
				for owner := 0; owner < players+1; owner++ {
					images[grid.Cell{
						Exploding:   exploding,
						Type:        typeMapping(simple),
						Orientation: grid.Orientation(orientation),
						Owner:       grid.Owner(owner)}] = surface
				}
			}
		}
	}
	for _, rotateable := range []string{"l", "p"} {
		for rotations := 0; rotations < 4; rotations++ {
			for player := 1; player <= players; player++ {
				for _, exploding := range []bool{false, true} {
					surface, err := loadImageRotated(fmt.Sprintf("final/images/%s%d.png",
						rotateable, (player-1)%2+1), float64(rotations)/4)
					if err != nil {
						cleanupImages(images)
						return nil, nil, err
					}
					if exploding {
						err = overlay(surface, "final/images/ex.png")
						if err != nil {
							surface.Free()
							cleanupImages(images)
							return nil, nil, err
						}
					}
					images[grid.Cell{
						Exploding:   exploding,
						Type:        typeMapping(rotateable),
						Owner:       grid.Owner(player),
						Orientation: grid.Orientation(rotations)}] = surface
				}
			}
		}
	}

	rotations = make(map[grid.Owner][]*sdl.Surface)
	for _, rotateable := range []string{"p"} {
		for player := 1; player <= players; player++ {
			for frame := 0; frame < 4*frames; frame++ {
				surface, err := loadImageRotated(fmt.Sprintf("final/images/%s%d.png",
					rotateable, (player-1)%2+1), float64(frame)/(4*float64(frames)))
				if err != nil {
					cleanupImages(images)
					cleanupRotations(rotations)
					return nil, nil, err
				}
				rotations[grid.Owner(player)] = append(
					rotations[grid.Owner(player)], surface)
			}
		}
	}

	return images, rotations, nil
}

func overlay(surface *sdl.Surface, img_name string) error {
	data, err := assets.Asset(img_name)
	if err != nil {
		return err
	}
	mem := sdl.RWFromMem(unsafe.Pointer(&data[0]), len(data))
	ontop, err := img.Load_RW(mem, 0)
	defer ontop.Free()
	if err != nil {
		return err
	}
	return ontop.BlitScaled(nil, surface, nil)
}

func typeMapping(name string) grid.Type {
	switch name {
	case "battery":
		return grid.Battery
	case "floor":
		return grid.Empty
	case "wall":
		return grid.Wall
	case "p":
		return grid.Player
	case "l":
		return grid.Laser
	default:
		panic("unknown type")
	}
}

func loadImageRotated(image_name string, angle float64) (*sdl.Surface, error) {
	file, err := assets.Asset(image_name)
	if err != nil {
		return nil, err
	}
	i, _, err := image.Decode(bytes.NewReader(file))
	if err != nil {
		return nil, err
	}

	dst := image.NewRGBA(i.Bounds())
	err = graphics.Rotate(dst, i, &graphics.RotateOptions{
		Angle: angle * tau})
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	err = png.Encode(&buf, dst)
	if err != nil {
		return nil, err
	}
	b := buf.Bytes()

	return img.Load_RW(sdl.RWFromMem(unsafe.Pointer(&b[0]), len(b)), 0)
}
