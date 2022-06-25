package display

import (
	"fmt"
	"github.com/jphastings/jan-poka/pkg/projections"
	"github.com/veandco/go-sdl2/sdl"
	"image"
	"image/color"
)

const (
	w = 1920
	h = 1080
)

type sdlDrawable struct {
	window  *sdl.Window
	surface *sdl.Surface
}

func (d sdlDrawable) Close() {
	d.window.Destroy()
	sdl.Quit()
}

func Get() (Drawable, error) {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return nil, err
	}

	window, err := sdl.CreateWindow(
		"Emulated Jan-Poka projector",
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		w, h, sdl.WINDOW_SHOWN)
	if err != nil {
		return nil, err
	}

	surface, err := window.GetSurface()
	if err != nil {
		return nil, err
	}
	_ = surface.FillRect(nil, 0)
	_ = window.UpdateSurface()

	return sdlDrawable{
		window:  window,
		surface: surface,
	}, nil
}

func (d sdlDrawable) Loop() {
	fmt.Println("Looping")
	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quitting UI")
				running = false
				break
			}
		}
	}
}

func (d sdlDrawable) Bounds() image.Rectangle {
	return d.surface.Bounds()
}

func (d sdlDrawable) ShowImage(i image.Image, pos projections.Pos) {
	_ = i
	_ = pos
	// TODO: Implement this
}

func (d sdlDrawable) DrawPoint(pos projections.Pos) {
	d.surface.Set(int(pos.X), int(pos.Y), color.RGBA{R: 255, A: 255})
	_ = d.window.UpdateSurface()
}
