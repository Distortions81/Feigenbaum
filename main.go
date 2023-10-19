package main

import (
	"image"
	"image/color"
	"math"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 1024
	screenHeight = 1024
	supersample  = 8.0
)

type Game struct {
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	imgLock.Lock()
	op := &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear}
	op.GeoM.Scale(1.0/supersample, 1.0/supersample)
	screen.DrawImage(imgBuf, op)
	imgLock.Unlock()
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Feigenbaum Constant")
	if err := ebiten.RunGameWithOptions(newGame(), &ebiten.RunGameOptions{GraphicsLibrary: ebiten.GraphicsLibraryOpenGL}); err != nil {
		return
	}
}

var imgBuf *ebiten.Image
var imgLock sync.Mutex

func newGame() *Game {

	op := &ebiten.NewImageOptions{Unmanaged: true}
	rect := image.Rectangle{}
	rect.Max.X = screenWidth * supersample
	rect.Max.Y = screenHeight * supersample
	imgBuf = ebiten.NewImageWithOptions(rect, op)
	imgBuf.Fill(color.White)
	go calc()

	return &Game{}
}

func calc() {
	time.Sleep(time.Second)
	start := 3.4

	for r := start; r < 4.0; r += 0.000002 {
		x := 0.5
		for i := 0; i < 100; i++ {
			x = r * x * (1.0 - x)
		}
		past := make([]float64, 1)
		past[0] = x
	DONE:
		for i := 0; i < 100; i++ {
			x = r * x * (1.0 - x)
			for _, n := range past {
				eps := math.Abs(x - n)
				if eps < math.SmallestNonzeroFloat64 {
					break DONE
				}
			}
			past = append(past, x)

		}
		for _, point := range past {
			x := (r - start) * (float64(screenWidth) * 1.6 * supersample)
			y := point * (float64(screenHeight) * supersample)

			imgLock.Lock()
			imgBuf.Set(int(x), int(y), color.Black)
			imgLock.Unlock()
		}
	}
}
