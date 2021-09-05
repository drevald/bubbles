package main

import (
	"fmt"
	"math/rand"
	"image/color"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Matrix struct {
	cells []int
	width int
	height int
}

func (m *Matrix) get(i, j int) int {
	return m.cells[j * m.width + i]
} 

func (m *Matrix) set(i, j, value int) {
	m.cells[j * m.width + i] = value
} 

func (m *Matrix) RotateRight() *Matrix {
	fmt.Println("rotating right")
	mRot := &Matrix{cells:make([]int, len(m.cells)), width:m.height, height: m.width }
	for i:=0; i < m.width; i++ {
		for j := 0; j < m.height; j++ {
			mRot.set(j, i, m.get(m.width - i - 1, j))
		}
	}
	return mRot;
}

func (m *Matrix) RotateLeft() *Matrix {
	return m.RotateRight().RotateRight().RotateRight()
}

type Game struct {
	counter int
	field *Matrix
	block *Matrix
	blockX int
	blockY int
	cell_size int
	freq int
	colors []color.RGBA	
}

func (g *Game) Drop() {
	g.freq = 1
}

var blocks = []Matrix{
	{cells:[]int{1, 1, 1, 1}, width: 2, height: 2},
	{cells:[]int{0, 2, 0, 2, 2, 2}, width: 3, height: 2},
	{cells:[]int{3, 0, 0, 3, 3, 3}, width: 3, height: 2},
	{cells:[]int{0, 0, 4, 4, 4, 4}, width: 3, height: 2},
	{cells:[]int{5, 5, 0, 0, 5, 5}, width: 3, height: 2},
	{cells:[]int{0, 6, 6, 6, 6, 0}, width: 3, height: 2},
	{cells:[]int{7, 7, 7, 7}, width: 1, height: 4},
}

func (g *Game) Draw (screen *ebiten.Image) {

	g.counter++
	if g.block == nil {
		fmt.Println(rand.Intn(6))
		g.block = &blocks[rand.Intn(7)]
		g.blockX = 0
		g.blockY = 0
	} else if (g.counter % g.freq == 0) {
		g.blockY += 1
	}
	
	image := ebiten.NewImage(g.cell_size, g.cell_size)
	options := &ebiten.DrawImageOptions {}
	screen.Fill(color.RGBA64{255, 0, 255, 255})
	image.Fill(g.colors[1])

	// Draw glass
	for i:=0; i < g.field.width; i++ {
		for j := 0; j < g.field.height; j++ {
			image.Fill(g.colors[g.field.get(i, j)])
			options.GeoM.Reset()
			options.GeoM.Translate(float64(i*g.cell_size), float64(j*g.cell_size))
			screen.DrawImage(image, options)
		}		
	}

	// Draw block	
	for i:=0; i<g.block.width; i++ {
		for j:=0; j<g.block.height; j++ {
			image.Fill(g.colors[g.block.get(i,j)])
			options.GeoM.Reset()
			options.GeoM.Translate(float64((g.blockX + i)*g.cell_size), float64((g.blockY + j)*g.cell_size))
			if (g.block.get(i, j) > 0) {
				screen.DrawImage(image, options)
			}			
		}
	}

	if (g.BlockLanded()) {
		g.MergeBlock()
		g.RemoveLines()
	}

}

func (g *Game) BlockLanded() bool {
	for i:=0; i<g.block.width; i++ {
		for j:=0; j<g.block.height; j++ {
			if (g.block.get(i, j) > 0 && j + g.blockY + 1>= g.field.height) {
				fmt.Println(g.block)
				fmt.Printf("%d + %d >= %d\n", j, g.blockY, g.field.height)
				return true
			}
			if (g.block.get(i, j) > 0 && g.field.get(i + g.blockX, j + g.blockY + 1) > 0) {
				return true
			}
		}
	}	
	return false
}

func (g *Game) MergeBlock() {
	for i:=0; i<g.block.width; i++ {
		for j:=0; j<g.block.height; j++ {
			fmt.Printf("g.field.set(%d + %d, %d + %d, %d)\n", i, g.blockX, j, g.blockY, g.block.get(i, j))
			if(g.block.get(i, j) > 0) {
				g.field.set(i + g.blockX, j + g.blockY, g.block.get(i, j))
			}			
		}
	}	
	g.block = nil
	g.freq = 25
}

func (g *Game) LineFull (j int) bool {
	for i:= 0; i < g.field.width; i++ {
		if g.field.get(i, j) == 0 {
			return false
		}
	}
	return true
}

func (g *Game) RemoveLine(j int) {
	for k :=j; k>1; k-- {
		for i :=0; i<g.field.width; i++ {
			g.field.set(i, k, g.field.get(i, k-1))
		}
	}
}

func (g *Game) RemoveLines() {
	for j:=g.field.height-1; j>0; j-- {
		for g.LineFull(j) {
			g.RemoveLine(j)
		}
	}
}

func (g *Game) Update () error {
	
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		g.blockX -= 1
	} 
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		g.blockX += 1
	} 
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		g.block = g.block.RotateRight()
	} 
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		g.block = g.block.RotateLeft()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.Drop()
	}	
	return nil
}

func (g *Game) Layout (outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {	
	return g.field.width * g.cell_size, g.field.height * g.cell_size
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
			{32, 32, 32, 255},
			{255, 0, 0, 255},
			{0, 255, 0, 255},
			{0, 0, 255, 255},
			{255, 255, 0, 255},
			{255, 0, 255, 255},
			{0, 255, 255, 255},
			{255, 255, 255, 255},			
		},
		freq:25,
	}) 
}