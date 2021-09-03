package main

import (
	"fmt"
	"image/color"
	"github.com/hajimehoshi/ebiten/v2"
)

type Matrix struct {
	cells []int
	width int
	height int
}

func (m *Matrix) get(x, y int) int {
	fmt.Println("Cell at " + x + ", " + y + " is " + (y * m.width + x))
	return m.cells[y * m.width + x]
} 

type Game struct {
	counter int
	field *Matrix
	block *Matrix
	blockX int
	blockY int
	cell_size int
	colors []color.RGBA	
}

func (g *Game) Draw (screen *ebiten.Image) {
	fmt.Println("Draw")
	g.counter++
	if g.block == nil {
		g.block = &Matrix{cells:make([]int, 4), width:2, height:2}
		g.blockX = 0
		g.blockY = 0
	} else if (g.counter % 10 == 0) {
		g.blockY += 1
	}
	
	image := ebiten.NewImage(g.cell_size, g.cell_size)
	options := &ebiten.DrawImageOptions {}

	for i:=0; i < g.field.width; i++ {
		for j := 0; j < g.field.height; j++ {
			image.Fill(g.colors[g.field.get(i, j)])
			options.GeoM.Translate(float64(i*g.cell_size), float64(j*g.cell_size))
			screen.DrawImage(image, options)
		}		
	}

	// imageOne := ebiten.NewImage(20, 20)
	// imageTwo := ebiten.NewImage(10, 10)
	// imageOne.Fill(color.RGBA{255, 255, 0, 255})
	// imageTwo.Fill(color.RGBA{0, 255, 0, 255})
	// optionsOne := &ebiten.DrawImageOptions {}
	// optionsOne.GeoM.Translate(float64((g.counter/2) % 100), 0)
	// screen.DrawImage(imageOne, optionsOne)
	// optionsTwo := &ebiten.DrawImageOptions {}
	// optionsTwo.GeoM.Translate(0, float64(g.counter % 100))
	// screen.DrawImage(imageTwo, optionsTwo)
}

func (g *Game) Update () error {
	fmt.Println("Update")
	return nil
}

func (g *Game) Layout (outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	fmt.Println("Layout")
	return g.field.width * g.cell_size + 200, g.field.height * g.cell_size
}

func main() {	
	fmt.Println("Start Game")
	ebiten.RunGame(&Game{
		field: &Matrix{
			cells:make([]int, 10*20),
			width:10,
			height:20,			
		},
		cell_size: 10,
		colors: []color.RGBA{
			{255, 0, 0, 255},
			{255, 0, 0, 255},
			{0, 255, 0, 255},
			{0, 0, 255, 255},
			{255, 255, 0, 255},
			{255, 0, 255, 255},
			{0, 255, 255, 255},
		},
	}) 
}
