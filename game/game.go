package game

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	_ "image/png"
	"math/rand"
	"bytes"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	_ "image/png"	
	"embed"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"

)

const (
	
	frameOX     = 0
	frameOY     = 0
	frameWidth  = 84
	frameHeight = 250
	frameCount  = 15
)

//go:embed bubble.png
//go:embed weed.png
var f embed.FS

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
	fmt.Println("Rotating right")
	mRot := &Matrix{cells:make([]int, len(m.cells)), width:m.height, height: m.width }
	for i:=0; i < m.width; i++ {
		for j := 0; j < m.height; j++ {
			mRot.set(j, i, m.get(m.width - i - 1, j))
		}
	}
	return mRot;
}

func (m *Matrix) RotateLeft() *Matrix {
	m.debug()
	fmt.Println("Rotating left")
	m1 := m.RotateRight().RotateRight().RotateRight()
	m.debug()
	return m1
}

func (m *Matrix) debug() {
	for j := 0; j < m.height; j++ {
		for i:=0; i < m.width; i++ {
			fmt.Print(m.get(i,j));
			fmt.Print("\t");
		}
		fmt.Print("\n");
	}
}

type Game struct {
	counter int
	field *Matrix
	block *Matrix
	blockX int
	blockY int
	cell_size int
	freq int
	over bool
	bubbleImage *ebiten.Image
	weedImage *ebiten.Image
	pressed bool
	screenWidth int
	screenHeight int
	zoom float64
}

func (g *Game) Init() {

	data, _ := f.ReadFile("bubble.png")
	bubbleImageReader := bytes.NewReader(data)
	bubbleImageDecoded, _ := png.Decode(bubbleImageReader)

	data, _ = f.ReadFile("weed.png")
	weedImageReader := bytes.NewReader(data)
	weedImageDecoded, _ := png.Decode(weedImageReader)

	g.field = &Matrix{
		cells:make([]int, 10*20),
		width:10,
		height:20,			
	 }
	
	g.zoom = 1.0
	g.cell_size = 10
	g.freq = 10
	g.over = false
	g.bubbleImage = ebiten.NewImageFromImage(bubbleImageDecoded)
	g.weedImage = ebiten.NewImageFromImage(weedImageDecoded)

	tt, _ := opentype.Parse(fonts.MPlus1pRegular_ttf)
	const dpi = 72

	mplusBigFont, _ = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    18,
		DPI:     dpi,// Use quantization to save glyph cache images.
	})

	normFont, _ = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    12,
		DPI:     dpi,// Use quantization to save glyph cache images.
	}) 
	
	g.pressed = false

}

func (g *Game) Drop() {
	g.freq = 1
}

var colors []color.RGBA	= []color.RGBA{	
	{111, 119, 89, 255},
	{75, 83, 32, 255},
	{5, 5, 5, 255},
	{255, 0, 0, 255},
	{0, 255, 0, 255},
	{0, 0, 255, 255},
	{255, 255, 0, 255},
	{255, 0, 255, 255},
	{0, 255, 255, 255},
	{255, 255, 255, 255},			
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

var (
	mplusBigFont    font.Face
	normFont		font.Face
)

func (g *Game) Draw (screen *ebiten.Image) {

	//fmt.Println("Width is " + screen.Bounds().String())

	g.counter++

	if g.block == nil && !g.over {

		g.block = &blocks[rand.Intn(7)]
		g.blockX = 0
 		g.blockY = g.field.height - g.block.height - 1 
		for i:=0; i<g.block.width; i++ {
			for j:=0; j<g.block.height; j++ {
				if (g.block.get(i, j) > 0 && g.field.get(g.blockX + i, g.blockY + j) > 0) {
					g.over = true
					return							
				}			
			}
		}	 

	} else if (g.counter % g.freq == 0 && !g.over) {
		g.blockY -= 1
	}
	
	options := &ebiten.DrawImageOptions {}
	screen.Fill(colors[0])

	options.GeoM.Translate(0, 60)
	options.GeoM.Scale(g.zoom, g.zoom)
	i := (g.counter / 10) % 15
	sx, sy := frameOX+i*frameWidth, frameOY
	screen.DrawImage(g.weedImage.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image), options)
	options.GeoM.Translate(20, 0)
	i = (i + 5) % 15
	sx, sy = frameOX+i*frameWidth, frameOY
	screen.DrawImage(g.weedImage.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image), options)

//	Draw glass
	for i:=0; i < g.field.width; i++ {
		for j := 0; j < g.field.height; j++ {
			options.GeoM.Reset()
			options.GeoM.Scale(g.zoom, g.zoom)
			options.GeoM.Translate(float64(i) * float64(g.cell_size) * g.zoom, float64(j) * float64(g.cell_size) * g.zoom)
			if (g.field.get(i,j) > 0) {
				screen.DrawImage(g.bubbleImage, options)
			} 			
		}		
	}

	// Draw block	
	if (!g.over) {

		for i:=0; i<g.block.width; i++ {
			for j:=0; j<g.block.height; j++ {
				options.GeoM.Reset()
				options.GeoM.Scale(g.zoom, g.zoom)
				options.GeoM.Translate(float64(g.blockX + i) * float64(g.cell_size) * g.zoom, float64(g.blockY + j) * float64(g.cell_size) * g.zoom)
				if (g.block.get(i, j) > 0) {
					screen.DrawImage(g.bubbleImage, options)
				}			
			}
		}	

		if (g.BlockLanded()) {
			g.MergeBlock()
			g.RemoveLines()
		}

	}

	if (g.over) {
		text.Draw(screen, "GAME\nOVER", mplusBigFont, 20, 60, colors[1])
		if (g.pressed) {
			vector.DrawFilledRect(screen, 22, 102, 60, 20, colors[1], false)
			text.Draw(screen, "RESTART", normFont, 27, 117, color.White)	
		} else {
			vector.DrawFilledRect(screen, 22, 102, 60, 20, color.White, false)
			vector.DrawFilledRect(screen, 20, 100, 60, 20, colors[1], false)
			text.Draw(screen, "RESTART", normFont, 25, 115, color.White)	
		}
	}

}

////////////////////////////////

func (g *Game) BlockLanded() bool {
	for i:=0; i<g.block.width; i++ {
		for j:=0; j<g.block.height; j++ {
			if (g.block.get(i, j) > 0 && j + g.blockY == 0) {
				return true
			}					
			if (g.block.get(i, j) > 0 && g.field.get(i + g.blockX, j + g.blockY - 1) > 0) {
				return true
			}
		}
	}	
	return false
}

func (g *Game) MergeBlock() {
	for i:=0; i<g.block.width; i++ {
		for j:=0; j<g.block.height; j++ {
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
	for k :=j; k < g.field.height-1; k++ {
		for i :=0; i<g.field.width; i++ {
			g.field.set(i, k, g.field.get(i, k+1))
		}
	}
}

func (g *Game) RemoveLines() {
	for j:=0; j < g.field.height; j++ {
		for g.LineFull(j) {
			g.RemoveLine(j)
		}
	}
}

func (g *Game) Update () error {

	//fmt.Println("Update")

	if (ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)) {
		x, y := ebiten.CursorPosition()
		if (x > 20 && x < 80 && y > 100 && y < 120) {
			g.pressed = true
		}
	} else if g.pressed {
		g.pressed = false
		g.over = false
		g.field = &Matrix{
			cells:make([]int, 10*20),
			width:10,
			height:20,			
		 }		
	} else {
		g.pressed = false
	}

	ids := inpututil.AppendJustPressedTouchIDs(nil)

	if ids != nil {
		for i := 0; i < len(ids); i++ {
			b := inpututil.IsTouchJustReleased(ids[i])
			z := inpututil.TouchPressDuration(ids[i])
			x, y := ebiten.TouchPosition(ids[i])
			fmt.Printf("Touch id = %d pos (%d, %d) dur %d released %t\n", ids[i], x, y, z, b)
			if (g.Intersects(g.block, x, y)) {
				g.Rotate()
			} else if (x < 50 && y > 50) {
				g.MoveLeft()
			} else if (x > 50 && y > 50) {
				g.MoveRight()				
			} else if (y < 50) {
				g.Drop()
			} else if (x > 20 && x < 80 && y > 100 && y < 120) {
					g.pressed = true
			} else if g.pressed {
				g.pressed = false
				g.over = false
				g.field = &Matrix{
					cells:make([]int, 10*20),
					width:10,
					height:20,			
				 }		
			} else {
				g.pressed = false
			}
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyI) {
		g.block.debug()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		g.MoveLeft();		
	} 
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		g.MoveRight()
	} 
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		g.UnRotate()		
	} 
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		g.Rotate()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.Drop()
	}	
	return nil
}

func (g *Game) MoveLeft() {
	g.blockX -= 1
	for i := 0; i < g.block.width; i++ {
		for j := 0; j < g.block.height; j++ {
			couldClashLeftSide := g.block.get(i, j) > 0 && (g.blockX + i) < 0
			couldClashLeftBlock := g.block.get(i, j) > 0 && g.field.get(g.blockX + i, g.blockY + j) > 0;
			if (couldClashLeftSide || couldClashLeftBlock) {
				g.blockX += 1
				break
			}
		}
	}
}

func (g *Game) MoveRight() {
	g.blockX += 1
	for i := 0; i < g.block.width; i++ {
		for j := 0; j < g.block.height; j++ {
			couldClashRightSide := g.block.get(i, j) > 0 && (g.blockX + i + 1) > g.field.width
			couldClashRightBlock := g.block.get(i, j) > 0 && g.field.get(g.blockX + i, g.blockY + j) > 0;
			if (couldClashRightSide || couldClashRightBlock) {
				g.blockX -= 1
				break
			}
		}
	}	
}

func (g *Game) Intersects (m *Matrix, x int, y int) bool {
	for i := 0; i < g.block.width; i++ {
		for j := 0; j < g.block.height; j++ {	
			if x > (g.blockX + i) * g.cell_size &&
			 x < (g.blockX + i + 1) * g.cell_size && 
			 y > (g.blockY + j) * g.cell_size && 
			 y < (g.blockY + j + 1) * g.cell_size && 
			 g.block.get(i, j) > 0 {
				return true
			}
		}
	}
	return false
}

func (g *Game) Rotate() {
	g.block = g.block.RotateRight()
	for i := 0; i < g.block.width; i++ {
		for j := 0; j < g.block.height; j++ {
			couldClashBlock := g.block.get(i, j) > 0 && g.field.get(g.blockX + i, g.blockY + j) > 0;
			couldClashRightSide := g.block.get(i, j) > 0 && (g.blockX + i + 1) > g.field.width
			couldClashLeftSide := g.block.get(i, j) > 0 && (g.blockX + i) < 0
			if (couldClashRightSide || couldClashBlock || couldClashLeftSide) {
				fmt.Println("Rotation failed")
				g.block = g.block.RotateLeft()
				break	
			}
		}
	}
}

func (g *Game) UnRotate() {
	g.block = g.block.RotateRight()
	for i := 0; i < g.block.width; i++ {
		for j := 0; j < g.block.height; j++ {
			couldClashBlock := g.block.get(i, j) > 0 && g.field.get(g.blockX + i, g.blockY + j) > 0;
			couldClashRightSide := g.block.get(i, j) > 0 && (g.blockX + i + 1) > g.field.width
			couldClashLeftSide := g.block.get(i, j) > 0 && (g.blockX + i) < 0
			if (!couldClashRightSide || !couldClashBlock || !couldClashLeftSide) {
				fmt.Println("Rotation failed")
				g.block = g.block.RotateLeft()
				break
			}
		}
	}
}

func (g *Game) Layout (outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {	
	//return g.field.width * g.cell_size, g.field.height * g.cell_size
	//fmt.Println("Layout")
	if (g.screenWidth == 0 || g.screenHeight == 0) {
		g.screenHeight = outsideHeight
		g.screenWidth = outsideWidth		
		g.zoom = float64(outsideHeight / 200)
		fmt.Printf("zoom is %f\n", float64(outsideHeight / 200))
	}
	return outsideWidth, outsideHeight
}
