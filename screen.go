package termloop

import "github.com/gdamore/tcell/v3"

// A Screen represents the current state of the display.
// To draw on the screen, create Drawables and set their positions.
// Then, add them to the Screen's Level, or to the Screen directly (e.g. a HUD).
type Screen struct {
	tcell.Screen
	oldCanvas Canvas
	canvas    Canvas
	level     Level
	Entities  []Drawable
	width     int
	height    int
	delta     float64
	fps       float64
	offsetx   int
	offsety   int
}

// NewScreen creates a new Screen, with no entities or level.
// Returns a pointer to the new Screen.
func NewScreen() *Screen {
	s, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}
	return NewScreenWithScreen(s)
}

func NewScreenWithScreen(s tcell.Screen) *Screen {
	sc := Screen{Screen: s, Entities: make([]Drawable, 0)}
	sc.canvas = NewCanvas(10, 10)
	return &sc
}

// Tick is used to process events such as input. It is called
// on every frame by the Game.
func (s *Screen) Tick(ev Event) {
	// TODO implement ticks using worker pools
	if s.level != nil {
		s.level.Tick(ev)
	}
	if ev != nil {
		for _, e := range s.Entities {
			e.Tick(ev)
		}
	}
}

// Draw is called every frame by the Game to render the current
// state of the screen.
func (s *Screen) Draw() {
	// Update termloop canvas
	s.canvas = NewCanvas(s.width, s.height)
	if s.level != nil {
		s.level.DrawBackground(s)
		s.level.Draw(s)
	}
	for _, e := range s.Entities {
		e.Draw(s)
	}
	// Check if anything changed between Draws
	if !s.canvas.equals(&s.oldCanvas) {
		// Draw to terminal
		termboxNormal(&s.canvas, s)
		s.Show()
	}
	s.oldCanvas = s.canvas
}

func (s *Screen) resize(w, h int) {
	s.width = w
	s.height = h
	c := NewCanvas(s.width, s.height)
	// Copy old data that fits
	for i := 0; i < min(s.width, len(s.canvas)); i++ {
		for j := 0; j < min(s.height, len(s.canvas[0])); j++ {
			c[i][j] = s.canvas[i][j]
		}
	}
	s.canvas = c
}

// SetLevel sets the Screen's current level to be l.
func (s *Screen) SetLevel(l Level) {
	s.level = l
}

// Level returns the Screen's current level.
func (s *Screen) Level() Level {
	return s.level
}

// AddEntity adds a Drawable to the current Screen, to be rendered.
func (s *Screen) AddEntity(d Drawable) {
	s.Entities = append(s.Entities, d)
}

// RemoveEntity removes Drawable d from the screen's entities.
func (s *Screen) RemoveEntity(d Drawable) {
	for i, elem := range s.Entities {
		if elem == d {
			s.Entities = append(s.Entities[:i], s.Entities[i+1:]...)
			return
		}
	}
}

// TimeDelta returns the number of seconds since the previous
// frame was rendered. Can be used for timings and animation.
func (s *Screen) TimeDelta() float64 {
	return s.delta
}

// Set the screen framerate.  By default, termloop will draw the
// the screen as fast as possible, which may use a lot of system
// resources.
func (s *Screen) SetFps(f float64) {
	s.fps = f
}

// RenderCell updates the Cell at a given position on the Screen
// with the attributes in Cell c.
func (s *Screen) RenderCell(x, y int, c *Cell) {
	newx := x + s.offsetx
	newy := y + s.offsety
	if newx >= 0 && newx < len(s.canvas) &&
		newy >= 0 && newy < len(s.canvas[0]) {
		renderCell(&s.canvas[newx][newy], c)
	}
}

func (s *Screen) offset() (int, int) {
	return s.offsetx, s.offsety
}

func (s *Screen) setOffset(x, y int) {
	s.offsetx, s.offsety = x, y
}

func renderCell(old, new_ *Cell) {
	if new_.Ch != 0 {
		old.Ch = new_.Ch
	}
	if new_.Bg != 0 {
		old.Bg = new_.Bg
	}
	if new_.Fg != 0 {
		old.Fg = new_.Fg
	}
}

func termboxNormal(canvas *Canvas, s tcell.Screen) {
	for i, col := range *canvas {
		for j, cell := range col {
			s.SetContent(i, j, cell.Ch,
				nil, tcell.StyleDefault.Foreground(tcell.Color(cell.Fg)).Background(tcell.Color(cell.Bg)))
		}
	}
}
