// Copyright (C) 2015 Space Monkey, Inc.

package sdl

import (
	"fmt"
	"sync"
	"time"
	"unsafe"

	"github.com/spacemonkeygo/spacelog"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_mixer"
	"github.com/veandco/go-sdl2/sdl_ttf"

	"sm/final/assets"
	"sm/final/grid"
	"sm/final/renderer"
)

var (
	logger = spacelog.GetLogger()
)

const (
	padding    = 50
	frameRate  = 24 // Hz
	statusSize = 40
	nickWidth  = 100
)

type SDLRenderer struct {
	Title            string
	mtx              sync.Mutex
	window           *sdl.Window
	font             *ttf.Font
	previous         [][]grid.Cell
	images           map[grid.Cell]*sdl.Surface
	sounds           map[string]*mix.Chunk
	frame_time       time.Duration
	frames           int
	player_rotations map[grid.Owner][]*sdl.Surface
	players          int
	closed           bool
}

var _ renderer.Renderer = (*SDLRenderer)(nil)

func NewRenderer(title string, width, height, players int,
	update_duration time.Duration) (rv *SDLRenderer, err error) {
	return rv, exec(func() error {
		window, err := sdl.CreateWindow(title, sdl.WINDOWPOS_UNDEFINED,
			sdl.WINDOWPOS_UNDEFINED, width, height, sdl.WINDOW_SHOWN)
		if err != nil {
			rv = nil
			return err
		}
		data, err := assets.Asset("final/Fixedsys500c.ttf")
		if err != nil {
			window.Destroy()
			rv = nil
			return err
		}
		mem := sdl.RWFromMem(unsafe.Pointer(&data[0]), len(data))
		font, err := ttf.OpenFontRW(mem, 0, 128)
		if err != nil {
			window.Destroy()
			rv = nil
			return err
		}
		frames := int(frameRate * update_duration / time.Second)
		images, player_rotations, err := loadImages(frames, players)
		if err != nil {
			window.Destroy()
			font.Close()
			rv = nil
			return err
		}
		sounds, err := loadSounds()
		if err != nil {
			cleanupImages(images)
			cleanupRotations(player_rotations)
			window.Destroy()
			font.Close()
			rv = nil
			return err
		}

		var frame_time time.Duration
		if frames > 0 {
			frame_time = update_duration / time.Duration(frames)
		}

		rv = &SDLRenderer{Title: title,
			window:           window,
			font:             font,
			images:           images,
			sounds:           sounds,
			frame_time:       frame_time,
			frames:           frames,
			player_rotations: player_rotations,
			closed:           false}
		return nil
	})
}

func (r *SDLRenderer) Close() {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	if r.closed {
		return
	}
	r.closed = true

	exec(func() error {
		cleanupImages(r.images)
		cleanupRotations(r.player_rotations)
		cleanupSounds(r.sounds)
		r.font.Close()
		r.window.Destroy()
		return nil
	})
}

func copyGrid(cells [][]grid.Cell) [][]grid.Cell {
	new_grid := make([][]grid.Cell, 0, len(cells))
	for _, row := range cells {
		new_row := make([]grid.Cell, 0, len(row))
		for _, cell := range row {
			new_row = append(new_row, cell)
		}
		new_grid = append(new_grid, new_row)
	}
	return new_grid
}

func (r *SDLRenderer) Clear() error {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	if r.closed {
		return nil
	}
	return exec(r.clear)
}

func (r *SDLRenderer) clear() error {
	r.previous = nil
	window, err := r.window.GetSurface()
	if err != nil {
		return err
	}
	return window.FillRect(nil, 0xff000000)
}

func (r *SDLRenderer) getBehind(m, n int, orientation grid.Orientation) (
	behind_m, behind_n int) {
	delta_x, delta_y := orientation.Delta()
	m -= delta_y
	n -= delta_x
	if m < 0 {
		m += len(r.previous)
	}
	if n < 0 {
		n += len(r.previous[0])
	}
	return m % len(r.previous), n % len(r.previous[0])
}

func (r *SDLRenderer) getAhead(m, n int, orientation grid.Orientation) (
	behind_m, behind_n int) {
	delta_x, delta_y := orientation.Delta()
	m += delta_y
	n += delta_x
	if m < 0 {
		m += len(r.previous)
	}
	if n < 0 {
		n += len(r.previous[0])
	}
	return m % len(r.previous), n % len(r.previous[0])
}

func (r *SDLRenderer) Update(cells [][]grid.Cell) error {
	var moves [][]int
	var rotations [][]int
	var shots [][]int
	explosions := map[grid.Type]bool{}
	var redraw, powerup bool
	var height, width, window_width, window_height int
	var cell_width, cell_height int32
	var surface *sdl.Surface

	r.mtx.Lock()
	defer r.mtx.Unlock()
	if r.closed {
		return nil
	}

	err := exec(func() (err error) {

		height = len(cells)
		if height < 1 {
			return r.clear()
		}
		width = len(cells[0])
		if width < 1 {
			return r.clear()
		}

		redraw = false
		if len(r.previous) != height || len(r.previous[0]) != width {
			r.previous = nil
		}
		if r.previous == nil || r.frames < 1 {
			redraw = true
		}

		window_width, window_height = r.window.GetSize()
		window_height -= statusSize
		cell_width = int32(window_width / width)
		cell_height = int32(window_height / height)

		surface, err = r.window.GetSurface()
		if err != nil {
			return err
		}

		for m, row := range cells {
			for n, cell := range row {
				should_redraw := true
				if r.previous != nil {
					should_redraw = func() bool {
						prev := r.previous[m][n]
						if cell.Exploding {
							// was there an explosion?
							explosions[prev.Type] = true
						}
						if prev.Type == grid.Battery && cell.Type == grid.Player {
							powerup = true
						}

						behind_m, behind_n := r.getBehind(m, n, cell.Orientation)
						behind := r.previous[behind_m][behind_n]
						switch {
						case cell == prev:
							// no change
							return false
						case cell.Owner == prev.Owner && cell.Type == prev.Type &&
							cell.Type == grid.Player && cell.Exploding == prev.Exploding:
							// player rotation
							rotations = append(rotations, []int{m, n})
							return false
						case cell.Type == behind.Type &&
							cell.Orientation == behind.Orientation &&
							cell.Owner == behind.Owner && (cell.Type == grid.Player ||
							cell.Type == grid.Laser):
							// normal player or laser move
							moves = append(moves, []int{behind_m, behind_n, m, n})
							return false
						case behind.Owner == cell.Owner &&
							behind.Orientation == cell.Orientation &&
							behind.Type == grid.Player && cell.Type == grid.Laser:
							// initial laser fire
							moves = append(moves, []int{behind_m, behind_n, m, n})
							shots = append(moves, []int{behind_m, behind_n})
							return false
						}
						return true
					}() || redraw
				}
				if should_redraw {
					err := r.drawCell(surface, r.getCellImage(cell),
						int32(n)*cell_width,
						int32(m)*cell_height,
						cell_width, cell_height, false)
					if err != nil {
						return err
					}
				}
			}
		}
		if sound, ok := r.sounds["laser1"]; len(shots) > 0 && ok {
			sound.PlayTimed(1, 0, 0)
		}
		if sound, ok := r.sounds["playerhit1"]; explosions[grid.Player] && ok {
			sound.PlayTimed(2, 0, 0)
		} else {
			if sound, ok := r.sounds["wallhit1"]; explosions[grid.Wall] && ok {
				sound.PlayTimed(3, 0, 0)
			}
		}
		if sound, ok := r.sounds["batteryhit1"]; (explosions[grid.Battery] ||
			explosions[grid.Laser] || explosions[grid.Empty]) && ok {
			sound.PlayTimed(4, 0, 0)
		}
		if sound, ok := r.sounds["batteryget1"]; powerup && ok {
			sound.PlayTimed(5, 0, 0)
		}

		return r.window.UpdateSurface()
	})
	if err != nil {
		return err
	}

	for frame := 1; frame <= r.frames; frame++ {
		if !redraw {
			err := exec(func() error {
				for _, mv := range moves {
					old_m, old_n, new_m, new_n := mv[0], mv[1], mv[2], mv[3]
					err := r.drawCell(surface, r.getCellImage(cells[old_m][old_n]),
						int32(old_n)*cell_width,
						int32(old_m)*cell_height,
						cell_width, cell_height, false)
					if err != nil {
						return err
					}
					err = r.drawCell(surface, r.getCellImage(r.previous[new_m][new_n]),
						int32(new_n)*cell_width,
						int32(new_m)*cell_height,
						cell_width, cell_height, false)
					if err != nil {
						return err
					}
				}

				for _, rot := range rotations {
					m, n := rot[0], rot[1]
					cell := cells[m][n]
					old_cell := r.previous[m][n]
					direction := 1
					if cell.Orientation.Right() == old_cell.Orientation {
						direction = -1
					}

					angle := int(old_cell.Orientation)*r.frames + direction*frame
					if angle < 0 {
						angle += r.frames * 4
					}
					angle %= r.frames * 4
					cell_surface := r.player_rotations[cell.Owner][angle]

					err = r.drawCell(surface, cell_surface,
						int32(n)*cell_width, int32(m)*cell_height, cell_width, cell_height,
						false)
					if err != nil {
						return err
					}
				}

				for _, mv := range moves {
					old_m, old_n, new_m, new_n := mv[0], mv[1], mv[2], mv[3]
					cell := cells[new_m][new_n]
					delta_x, delta_y := cell.Orientation.Delta()
					new_x := int32(old_n)*cell_width +
						cell_width*int32(frame)*int32(delta_x)/int32(r.frames)
					new_y := int32(old_m)*cell_height +
						cell_height*int32(frame)*int32(delta_y)/int32(r.frames)
					err = r.drawCell(surface, r.getCellImage(cell), new_x, new_y,
						cell_width, cell_height, true)
					if err != nil {
						return err
					}
				}
				for _, shot := range shots {
					m, n := shot[0], shot[1]
					err := r.drawCell(surface, r.getCellImage(cells[m][n]),
						int32(n)*cell_width,
						int32(m)*cell_height,
						cell_width, cell_height, false)
					if err != nil {
						return err
					}
				}

				return r.window.UpdateSurface()
			})
			if err != nil {
				return err
			}
		}
		time.Sleep(r.frame_time)
	}

	return exec(func() error {
		r.previous = copyGrid(cells)
		if !redraw {
			for m, row := range cells {
				for n, cell := range row {
					err := r.drawCell(surface, r.getCellImage(cell),
						int32(n)*cell_width, int32(m)*cell_height, cell_width, cell_height, false)
					if err != nil {
						return err
					}
				}
			}
		}
		return r.window.UpdateSurface()
	})
}

func (r *SDLRenderer) getCellImage(cell grid.Cell) *sdl.Surface {
	rv := r.images[cell]
	if rv == nil {
		logger.Critf("unknown cell type: %#v", cell)
	}
	return rv
}

func (r *SDLRenderer) drawCell(window_surf, cell_surface *sdl.Surface,
	x, y, w, h int32, transparent bool) error {
	if cell_surface == nil {
		return fmt.Errorf("unknown cell type")
	}
	if !transparent {
		err := r.getCellImage(grid.Cell{Type: grid.Empty}).BlitScaled(nil,
			window_surf, &sdl.Rect{X: x, Y: y + statusSize, W: w, H: h})
		if err != nil {
			return err
		}
	}
	return cell_surface.BlitScaled(nil, window_surf, &sdl.Rect{
		X: x, Y: y + statusSize, W: w, H: h})
}

func (r *SDLRenderer) textToSurface(msg string, dest *sdl.Surface,
	dest_rect *sdl.Rect) error {
	text, err := r.font.RenderUTF8_Solid(msg,
		sdl.Color{R: 0xff, G: 0xff, B: 0xff, A: 0xff})
	if err != nil {
		return err
	}
	defer text.Free()

	text_reformatted, err := text.Convert(dest.Format, 0)
	if err != nil {
		return err
	}
	defer text_reformatted.Free()

	text_aspect := float64(text_reformatted.W) /
		float64(text_reformatted.H)
	dest_aspect := float64(dest_rect.W) / float64(dest_rect.H)

	dest_rect_copy := *dest_rect

	if text_aspect > dest_aspect {
		// text is wider
		dest_rect_copy.H = int32(float64(dest_rect.W) / text_aspect)
		dest_rect_copy.Y += (dest_rect.H - dest_rect_copy.H) / 2
	} else {
		// text is taller
		dest_rect_copy.W = int32(float64(dest_rect.H) * text_aspect)
		dest_rect_copy.X += (dest_rect.W - dest_rect_copy.W) / 2
	}

	return text_reformatted.BlitScaled(nil, dest, &dest_rect_copy)
}

func (r *SDLRenderer) Message(msg string, msgtype renderer.MessageType) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	if r.closed {
		return nil
	}

	return exec(func() error {
		if sound, ok := r.sounds["gong1"]; msgtype == renderer.GameStart && ok {
			sound.PlayTimed(6, 0, 0)
		}
		if sound, ok := r.sounds["applause1"]; msgtype == renderer.GameOver && ok {
			sound.PlayTimed(7, 0, 0)
		}

		window, err := r.window.GetSurface()
		if err != nil {
			return err
		}

		if r.previous == nil {
			err := r.clear()
			if err != nil {
				return err
			}

		} else {
			window_width, window_height := r.window.GetSize()
			window_height -= statusSize
			cell_width := int32(window_width / len(r.previous[0]))
			cell_height := int32(window_height / len(r.previous))
			for m, row := range r.previous {
				for n, cell := range row {
					err := r.drawCell(window, r.getCellImage(cell),
						int32(n)*cell_width, int32(m)*cell_height, cell_width, cell_height,
						false)
					if err != nil {
						return err
					}
				}
			}
		}

		window_sized := &sdl.Rect{
			X: padding,
			Y: padding,
			W: window.W - (padding * 2),
			H: window.H - (padding * 2)}

		err = r.textToSurface(msg, window, window_sized)
		if err != nil {
			return err
		}

		return r.window.UpdateSurface()
	})
}

func (r *SDLRenderer) WaitForQuit() {
	for {
		var window_closed bool
		exec(func() error {
			for {
				ev := sdl.PollEvent()
				if ev == nil {
					return nil
				}
				if wev, ok := ev.(*sdl.WindowEvent); ok &&
					wev.Event == sdl.WINDOWEVENT_CLOSE &&
					wev.WindowID == r.window.GetID() {
					window_closed = true
				}
			}
		})
		if window_closed {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func (r *SDLRenderer) SetStatus(statuses []renderer.PlayerStatus) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	if r.closed {
		return nil
	}

	return exec(func() error {
		for _, status := range statuses {
			if status.Health < 0 || status.Health > 1 {
				return fmt.Errorf("health should be between 0 and 1")
			}
			if status.Energy < 0 || status.Energy > 1 {
				return fmt.Errorf("energy should be between 0 and 1")
			}
		}

		window, err := r.window.GetSurface()
		if err != nil {
			return err
		}

		err = window.FillRect(&sdl.Rect{
			W: window.W, H: statusSize}, 0xff000000)
		if err != nil {
			return err
		}

		bar_width := window.W/2 - nickWidth

		if len(statuses) > 0 {
			p1 := statuses[0]
			err = r.textToSurface(p1.Moniker, window, &sdl.Rect{
				W: nickWidth, H: statusSize})
			if err != nil {
				logger.Crit("here")
				return err
			}
			err = window.FillRect(&sdl.Rect{
				X: nickWidth,
				W: int32(p1.Health * float64(bar_width)),
				H: 2 * statusSize / 3}, 0xff06d0ff)
			if err != nil {
				logger.Crit("here")
				return err
			}
			err = window.FillRect(&sdl.Rect{
				X: nickWidth,
				Y: 2 * statusSize / 3,
				W: int32(p1.Energy * float64(bar_width)),
				H: statusSize / 3}, 0xff17de9f)
			if err != nil {
				logger.Crit("here")
				return err
			}
		}

		if len(statuses) > 1 {
			p2 := statuses[1]
			err = r.textToSurface(p2.Moniker, window, &sdl.Rect{
				X: window.W - nickWidth, W: nickWidth, H: statusSize})
			if err != nil {
				logger.Crit("here")
				return err
			}
			p2health_width := int32(p2.Health * float64(bar_width))
			p2energy_width := int32(p2.Energy * float64(bar_width))
			err = window.FillRect(&sdl.Rect{
				X: window.W - nickWidth - p2health_width,
				W: p2health_width,
				H: 2 * statusSize / 3}, 0xffff5a12)
			if err != nil {
				logger.Crit("here")
				return err
			}
			err = window.FillRect(&sdl.Rect{
				Y: 2 * statusSize / 3,
				X: window.W - nickWidth - p2energy_width,
				W: p2energy_width,
				H: statusSize / 3}, 0xff17de9f)
			if err != nil {
				logger.Crit("here")
				return err
			}
		}

		return nil
	})
}
