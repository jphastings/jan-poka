package display

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/gonutz/framebuffer"
)

type fbDrawable struct {
	fb draw.Image
}

func Get() (Drawable, error) {
	fb, err := framebuffer.Open(rpi.HDMI)
	if err != nil {
		return nil, fmt.Errorf("could not connect to the framebuffer: %w", err)
	}
	return fbDrawable{fb: fb}, nil
}

func (d sdlDrawable) Loop() {
	exitSignal := make(chan os.Signal)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal
}

func (d sdlDrawable) Bounds() image.Rectangle {
	return d.Bounds()
}

func (d sdlDrawable) ShowImage(i image.Image, pos projections.Pos) {
	draw.Draw(d.fb, i.Bounds(), i, image.Point{X: int(pos.X), Y: int(pos.Y)}, draw.Src)
}

func (d sdlDrawable) DrawPoint(pos projections.Pos) {
	d.fb.(draw.Image).Set(int(pos.X), int(pox.Y), color.RGBA{R: 255, A: 255})
}
