package termloop

import (
	"strings"

	"github.com/gdamore/tcell/v3"
)

// A Canvas is a 2D array of Cells, used for drawing.
// The structure of a Canvas is an array of columns.
// This is so it can be addressed canvas[x][y].
type Canvas [][]Cell

// NewCanvas returns a new Canvas, with
// width and height defined by arguments.
func NewCanvas(width, height int) Canvas {
	canvas := make(Canvas, width)
	for i := range canvas {
		canvas[i] = make([]Cell, height)
	}
	return canvas
}

func (canvas *Canvas) equals(oldCanvas *Canvas) bool {
	c := *canvas
	c2 := *oldCanvas
	if c2 == nil {
		return false
	}
	if len(c) != len(c2) {
		return false
	}
	if len(c[0]) != len(c2[0]) {
		return false
	}
	for i := range c {
		for j := range c[i] {
			equal := c[i][j].equals(&(c2[i][j]))
			if !equal {
				return false
			}
		}
	}
	return true
}

// CanvasFromString returns a new Canvas, built from
// the characters in the string str. Newline characters in
// the string are interpreted as a new Canvas row.
func CanvasFromString(str string) Canvas {
	lines := strings.Split(str, "\n")
	runes := make([][]rune, len(lines))
	width := 0
	for i := range lines {
		runes[i] = []rune(lines[i])
		width = max(width, len(runes[i]))
	}
	height := len(runes)
	canvas := make(Canvas, width)
	for i := 0; i < width; i++ {
		canvas[i] = make([]Cell, height)
		for j := 0; j < height; j++ {
			if i < len(runes[j]) {
				canvas[i][j] = Cell{Ch: runes[j][i]}
			}
		}
	}
	return canvas
}

// Drawable represents something that can be drawn, and placed in a Level.
type Drawable interface {
	Tick(Event)   // Method for processing events, e.g. input
	Draw(*Screen) // Method for drawing to the screen
}

// Physical represents something that can collide with another
// Physical, but cannot process its own collisions.
// Optional addition to Drawable.
type Physical interface {
	Position() (int, int) // Return position, x and y
	Size() (int, int)     // Return width and height
}

// DynamicPhysical represents something that can process its own collisions.
// Implementing this is an optional addition to Drawable.
type DynamicPhysical interface {
	Position() (int, int) // Return position, x and y
	Size() (int, int)     // Return width and height
	Collide(Physical)     // Handle collisions with another Physical
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Abstract Termbox stuff for convenience - users
// should only need Termloop imported

// Represents a character to be drawn on the screen.
type Cell struct {
	Fg Attr // Foreground colour
	Bg Attr // Background color
	Ch rune // The character to draw
}

func (c *Cell) equals(c2 *Cell) bool {
	return c.Fg == c2.Fg &&
		c.Bg == c2.Bg &&
		c.Ch == c2.Ch
}

// Provides an event, for input, errors or resizing.
// Resizing and errors are largely handled by Termloop itself
// - this would largely be used for input.
type Event tcell.Event

func convertEvent(ev tcell.Event) Event {
	return Event(ev)
}

type (
	Attr      tcell.Color
	Key       = tcell.Key
	Modifier  = tcell.ModMask
	EventType = tcell.Event
)

// Types of event. For example, a keyboard press will be EventKey.
const (
	EventKey       = tcell.KeyEnter
	EventResize    = tcell.KeyRune
	EventMouse     = tcell.KeyLeft
	EventError     = tcell.KeyRight
	EventInterrupt = tcell.KeyDelete
	EventRaw       = tcell.KeyInsert
	EventNone      = tcell.KeyNUL
)

// Cell colors. You can combine these with multiple attributes using
// a bitwise OR ('|'). Colors can't combine with other colors.
const (
	ColorDefault = Attr(tcell.ColorDefault)
	ColorBlack   = Attr(tcell.ColorBlack)
	ColorRed     = Attr(tcell.ColorRed)
	ColorGreen   = Attr(tcell.ColorGreen)
	ColorYellow  = Attr(tcell.ColorYellow)
	ColorBlue    = Attr(tcell.ColorBlue)
	ColorMagenta = Attr(tcell.ColorPurple)
	ColorCyan    = Attr(tcell.ColorTeal)
	ColorWhite   = Attr(tcell.ColorWhite)
)

// Cell attributes. These can be combined with OR.
const (
	AttrBold    = tcell.AttrBold
	AttrReverse = tcell.AttrReverse
)

const ModAlt = tcell.ModAlt

// Key constants. See Event.Key.
const (
	KeyF1         = tcell.KeyF1
	KeyF2         = tcell.KeyF2
	KeyF3         = tcell.KeyF3
	KeyF4         = tcell.KeyF4
	KeyF5         = tcell.KeyF5
	KeyF6         = tcell.KeyF6
	KeyF7         = tcell.KeyF7
	KeyF8         = tcell.KeyF8
	KeyF9         = tcell.KeyF9
	KeyF10        = tcell.KeyF10
	KeyF11        = tcell.KeyF11
	KeyF12        = tcell.KeyF12
	KeyInsert     = tcell.KeyInsert
	KeyDelete     = tcell.KeyDelete
	KeyHome       = tcell.KeyHome
	KeyEnd        = tcell.KeyEnd
	KeyPgUp       = tcell.KeyPgUp
	KeyPgDn       = tcell.KeyPgDn
	KeyArrowUp    = tcell.KeyUp
	KeyArrowDown  = tcell.KeyDown
	KeyArrowLeft  = tcell.KeyLeft
	KeyArrowRight = tcell.KeyRight
	KeyCtrlA      = tcell.KeyCtrlA
	KeyCtrlB      = tcell.KeyCtrlB
	KeyCtrlC      = tcell.KeyCtrlC
	KeyCtrlD      = tcell.KeyCtrlD
	KeyCtrlE      = tcell.KeyCtrlE
	KeyCtrlF      = tcell.KeyCtrlF
	KeyCtrlG      = tcell.KeyCtrlG
	KeyBackspace  = tcell.KeyBackspace
	KeyCtrlH      = tcell.KeyCtrlH
	KeyTab        = tcell.KeyTab
	KeyCtrlI      = tcell.KeyCtrlI
	KeyCtrlJ      = tcell.KeyCtrlJ
	KeyCtrlK      = tcell.KeyCtrlK
	KeyCtrlL      = tcell.KeyCtrlL
	KeyEnter      = tcell.KeyEnter
	KeyCtrlM      = tcell.KeyCtrlM
	KeyCtrlN      = tcell.KeyCtrlN
	KeyCtrlO      = tcell.KeyCtrlO
	KeyCtrlP      = tcell.KeyCtrlP
	KeyCtrlQ      = tcell.KeyCtrlQ
	KeyCtrlR      = tcell.KeyCtrlR
	KeyCtrlS      = tcell.KeyCtrlS
	KeyCtrlT      = tcell.KeyCtrlT
	KeyCtrlU      = tcell.KeyCtrlU
	KeyCtrlV      = tcell.KeyCtrlV
	KeyCtrlW      = tcell.KeyCtrlW
	KeyCtrlX      = tcell.KeyCtrlX
	KeyCtrlY      = tcell.KeyCtrlY
	KeyCtrlZ      = tcell.KeyCtrlZ
	KeyEsc        = tcell.KeyEsc
	KeySpace      = tcell.KeyRune
	KeyBackspace2 = tcell.KeyBackspace2
)
