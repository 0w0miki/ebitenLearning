package game

import (
	"log"

	_ "image/jpeg"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type background struct {
	image  *ebiten.Image
	width  int
	height int
	y1     float64
	y2     float64
	speed  float64
}

func loadBackground(path string) *background {
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatal(err)
	}

	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	return &background{
		image:  img,
		width:  width,
		height: height,
		y1:     0,
		y2:     -float64(height),
		speed:  2,
	}
}

func (bg *background) draw(screen *ebiten.Image, cfg *config) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, bg.y1)
	screen.DrawImage(bg.image, op)
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, bg.y2)
	screen.DrawImage(bg.image, op)
}

func (bg *background) update() {
	bg.y1 += bg.speed
	bg.y2 += bg.speed
	if bg.y1 >= float64(bg.height) {
		bg.y1 = -float64(bg.height)
	}
	if bg.y2 >= float64(bg.height) {
		bg.y2 = -float64(bg.height)
	}
}
