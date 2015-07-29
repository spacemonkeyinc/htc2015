// Copyright (C) 2015 Space Monkey, Inc.

package grid

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spacemonkeygo/errors"
)

type Type string

var (
	GridError = errors.NewClass("grid error")
)

const (
	Empty   Type = "empty"
	Wall    Type = "wall"
	Player  Type = "player"
	Battery Type = "battery"
	Laser   Type = "laser"
)

func (t Type) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(t))
}

func (t *Type) UnmarshalJSON(p []byte) (err error) {
	var raw string
	err = json.Unmarshal(p, &raw)
	if err != nil {
		return err
	}
	switch Type(strings.ToLower(raw)) {
	case Empty:
		*t = Empty
	case Wall:
		*t = Wall
	case Player:
		*t = Player
	case Battery:
		*t = Battery
	case Laser:
		*t = Laser
	default:
		return errors.New(fmt.Sprintf("%s is not a valid cell type", raw))
	}
	return nil
}

type Orientation int

const (
	North Orientation = 0
	East  Orientation = 1
	South Orientation = 2
	West  Orientation = 3
)

func (o Orientation) String() string {
	switch o {
	case North:
		return "north"
	case East:
		return "east"
	case South:
		return "south"
	case West:
		return "west"
	}
	panic("unreachable")
}

func (o Orientation) Opposite() Orientation {
	switch o {
	case North:
		return South
	case East:
		return West
	case South:
		return North
	case West:
		return East
	}
	panic("unreachable")
}

func (o Orientation) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.String())
}

func (o *Orientation) UnmarshalJSON(p []byte) (err error) {
	var raw string
	err = json.Unmarshal(p, &raw)
	if err != nil {
		return err
	}
	switch strings.ToLower(raw) {
	case "north", "n":
		*o = North
	case "east", "e":
		*o = East
	case "south", "s":
		*o = South
	case "west", "w":
		*o = West
	default:
		// hmm, try as a number
		i, err := strconv.Atoi(raw)
		if err != nil {
			return err
		}
		switch Orientation(i) {
		case North:
			*o = North
		case East:
			*o = East
		case South:
			*o = South
		case West:
			*o = West
		default:
			return errors.New(fmt.Sprintf("%s is not a valid orientation", raw))
		}
	}
	return nil
}

func (o Orientation) Delta() (x, y int) {
	switch o {
	case North:
		return 0, -1
	case East:
		return 1, 0
	case South:
		return 0, 1
	case West:
		return -1, 0
	}
	panic("unreached")
}

func (o Orientation) Left() Orientation {
	switch o {
	case North:
		return West
	case East:
		return North
	case South:
		return East
	case West:
		return South
	}
	panic("unreached")
}

func (o Orientation) Right() Orientation {
	switch o {
	case North:
		return East
	case East:
		return South
	case South:
		return West
	case West:
		return North
	}
	panic("unreached")
}

func (o *Orientation) RotateLeft() {
	*o = o.Left()
}

func (o *Orientation) RotateRight() {
	*o = o.Right()
}

type Owner int

const (
	None Owner = 0
)

func (o Owner) String() string {
	switch o {
	case None:
		return "None"
	default:
		return fmt.Sprintf("Player%d", o)
	}
	panic("unreached")
}

type Cell struct {
	Type
	Orientation
	Owner
	Exploding bool
}

var (
	EmptyCell = Cell{Type: Empty}
	WallCell  = Cell{Type: Wall}
)

type Coord struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func (c Coord) String() string {
	return fmt.Sprintf("(%d,%d)", c.X, c.Y)
}

type Grid struct {
	cells [][]Cell
}

func NewEmpty(width, height int) *Grid {
	return &Grid{
		cells: newCells(width, height),
	}
}

func newCells(width, height int) (rv [][]Cell) {
	rv = make([][]Cell, 0, height)
	for m := 0; m < height; m++ {
		row := make([]Cell, 0, width)
		for n := 0; n < width; n++ {
			row = append(row, EmptyCell)
		}
		rv = append(rv, row)
	}
	return rv
}

func NewRandom(width, height int, walls int, enclosed bool) *Grid {
	r := rand.New(rand.NewSource(time.Now().Unix()))

	grid := NewEmpty(width, height)
	rv := grid.cells
	if enclosed {
		for n := 0; n < width; n++ {
			rv[0][n] = WallCell
			rv[height-1][n] = WallCell
		}
		for m := 1; m < height-1; m++ {
			rv[m][0] = WallCell
			rv[m][width-1] = WallCell
		}
	}
	for w := 0; w < walls; w++ {
		if r.Intn(2) == 0 {
			// horizontal
			length := r.Intn(width / 2)
			start := r.Intn(width - length)
			row := r.Intn(height)
			for i := 0; i < length; i++ {
				rv[row][start+i] = WallCell
			}
		} else {
			// vertical
			length := r.Intn(height / 2)
			start := r.Intn(height - length)
			col := r.Intn(width)
			for i := 0; i < length; i++ {
				rv[start+i][col] = WallCell
			}
		}
	}
	return grid
}

func (g *Grid) CopyTo(other *Grid) {
	// make sure height and width match, otherwise, we'll need to allocate
	// a new row/cols
	dest := other.cells
	if g.Width() != other.Width() || g.Height() != other.Height() {
		dest = newCells(g.Width(), g.Height())
	}

	for r, row := range g.cells {
		for c, col := range row {
			dest[r][c] = col
		}
	}

	other.cells = dest
}

func (g *Grid) Width() int {
	if g.Height() == 0 {
		return 0
	}
	return len(g.cells[0])
}

func (g *Grid) Height() int {
	return len(g.cells)
}

func (g *Grid) ClearCell(coord Coord) {
	cell := g.cellAt(coord)
	if cell != nil {
		*cell = EmptyCell
	}
}

func (g *Grid) SetCell(coord Coord, cell Cell) {
	mutable_cell := g.cellAt(coord)
	if mutable_cell != nil {
		*mutable_cell = cell
	}
}

func (g *Grid) SetCellExploding(coord Coord, exploding bool) {
	mutable_cell := g.cellAt(coord)
	if mutable_cell != nil {
		mutable_cell.Exploding = exploding
	}
}

func (g *Grid) Cells() [][]Cell {
	return g.cells
}

func (g *Grid) CellAt(coord Coord) Cell {
	if cell := g.cellAt(coord); cell != nil {
		return *cell
	}
	panic(fmt.Sprintf("no cell at %s", coord))
}

func (g *Grid) cellAt(coord Coord) *Cell {
	if coord.X >= g.Width() || coord.Y >= g.Height() {
		return nil
	}
	return &g.cells[coord.Y][coord.X]
}

func (g *Grid) RelativeTo(coord Coord, orientation Orientation) (
	rv Coord) {
	rv = coord
	dx, dy := orientation.Delta()
	rv.X += dx
	rv.Y += dy
	if rv.Y < 0 {
		rv.Y = g.Height() - 1
	}
	if rv.Y >= g.Height() {
		rv.Y = 0
	}
	if rv.X < 0 {
		rv.X = g.Width() - 1
	}
	if rv.X >= g.Width() {
		rv.X = 0
	}
	return rv
}

func (g *Grid) CellRelativeTo(coord Coord, orientation Orientation) (Cell,
	Coord) {
	new_coord := g.RelativeTo(coord, orientation)
	cell := g.cellAt(new_coord)
	if cell == nil {
		panic(fmt.Sprintf("expected cell at %s; got nil", coord))
	}
	return *cell, new_coord
}

func LoadFromFile(path string) (*Grid, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, GridError.Wrap(err)
	}
	defer file.Close()

	var cells [][]Cell
	var width int

	scanner := bufio.NewScanner(file)
	lineno := 0
	for ; scanner.Scan(); lineno++ {
		line := scanner.Text()
		if line == "" {
			continue
		}
		if width > 0 {
			if len(line) != width {
				return nil, GridError.New(
					"expected width %d on line %d, got %d",
					width, lineno, len(line))
			}
		} else {
			width = len(line)
		}
		row := make([]Cell, 0, len(line))
		for _, c := range line {
			switch c {
			case '_':
				row = append(row, EmptyCell)
			case 'W':
				row = append(row, WallCell)
			default:
				return nil, GridError.New("unexpected character %v on line %d",
					c, lineno)
			}
		}
		cells = append(cells, row)
	}
	if err = scanner.Err(); err != nil {
		return nil, GridError.Wrap(err)
	}
	if len(cells) == 0 {
		return nil, GridError.New("empty grid file")
	}
	return &Grid{
		cells: cells,
	}, nil
}

func (g *Grid) SerializeFor(owner Owner) string {
	var buf bytes.Buffer
	for y := 0; y < len(g.cells); y++ {
		for x := 0; x < len(g.cells[y]); x++ {
			cell := g.cells[y][x]

			r := '_'
			switch cell.Type {
			case Wall:
				r = 'W'
			case Player:
				if cell.Owner == owner {
					r = 'X'
				} else {
					r = 'O'
				}
			case Battery:
				r = 'B'
			case Laser:
				r = 'L'
			}
			buf.WriteRune(r)
		}
		buf.WriteRune('\n')
	}
	return buf.String()
}
