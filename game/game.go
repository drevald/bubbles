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
	frameWidth  = 85
	frameHeight = 250
	frameCount  = 15
	bottomHeight = 1
)

//go:embed bubble.png
//go:embed weed.png
//go:embed weed2.png
//go:embed bg.png
//go:embed water.png
//go:embed weed4.png
//go:embed hand.png
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
	lines int
	level int
	field *Matrix
	block *Matrix
	blockX int
	blockY int
	cell_size int
	freq int
	speed int
	over bool
	splash bool
	bubbleImage *ebiten.Image
	weedImage *ebiten.Image
	smallWeedImage *ebiten.Image
	bgImage *ebiten.Image
	waterImage *ebiten.Image
	handImage *ebiten.Image
	pressed bool
	splashPressed bool
	zoom float64
}

func (g *Game) Init() {

	data, _ := f.ReadFile("bubble.png")
	bubbleImageReader := bytes.NewReader(data)
	bubbleImageDecoded, _ := png.Decode(bubbleImageReader)

	_, _ = f.ReadFile("weed.png")

	data, _ = f.ReadFile("weed2.png")
	weedImageReader := bytes.NewReader(data)
	weedImageDecoded, _ := png.Decode(weedImageReader)

	data, _ = f.ReadFile("bg.png")
	bgImageReader := bytes.NewReader(data)
	bgImageDecoded, _ := png.Decode(bgImageReader)

	data, _ = f.ReadFile("water.png")
	waterImageReader := bytes.NewReader(data)
	waterImageDecoded, _ := png.Decode(waterImageReader)

	data, _ = f.ReadFile("weed4.png")
	smallWeedImageReader := bytes.NewReader(data)
	smallWeedImageDecoded, _ := png.Decode(smallWeedImageReader)

	data, _ = f.ReadFile("hand.png")
	handImageReader := bytes.NewReader(data)
	handImageDecoded, _ := png.Decode(handImageReader)

	g.field = &Matrix{
		cells:make([]int, 10*20),
		width:10,
		height:20,			
	 }
	
	g.lines = 0
	g.level = 1
	g.zoom = 1.0
	g.cell_size = 10
	g.speed = 15
	g.freq = 15
	g.over = false
	g.splash = true
	g.bubbleImage = ebiten.NewImageFromImage(bubbleImageDecoded)
	g.weedImage = ebiten.NewImageFromImage(weedImageDecoded)
	g.bgImage = ebiten.NewImageFromImage(bgImageDecoded)
	g.waterImage = ebiten.NewImageFromImage(waterImageDecoded)
	g.smallWeedImage = ebiten.NewImageFromImage(smallWeedImageDecoded)
	g.handImage = ebiten.NewImageFromImage(handImageDecoded)

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
	
	smallFont, _ = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    9,
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
	{34, 142, 143, 255},	
	{255, 79, 205, 255},	
}

var magenta color.RGBA = color.RGBA{255, 79, 205, 255};
var darkMagenta color.RGBA = color.RGBA{225, 0, 183, 255};

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
	smallFont		font.Face
)

func (g *Game) DrawBlock(screen *ebiten.Image, options *ebiten.DrawImageOptions) {

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

}

func (g *Game) Draw (screen *ebiten.Image) {

	//fmt.Println("Width is " + screen.Bounds().String())

	g.counter++

	if g.splash {
		if (g.counter > 400) {
			g.counter = 0
		} else if (g.counter > 0 && g.counter < 100) {
			bgOptions := &ebiten.DrawImageOptions {}
			screen.DrawImage(g.bgImage, bgOptions)
			screen.DrawImage(g.waterImage, bgOptions)
			g.block = &blocks[1]
			g.blockX = (g.field.width/2) - 1
			g.blockY = g.field.height/2 - 2
			g.blockX += (g.counter % 50)/10
			g.DrawBlock(screen, bgOptions)
			handOptions := &ebiten.DrawImageOptions {}
			handOptions.GeoM.Translate(70, 90);
			if (g.counter % 20 > 10) {
				screen.DrawImage(g.handImage, handOptions)			
			}		
			text.Draw(screen, "To move \nbubbbles right \ntap right side", smallFont, 10, 140, magenta);
		} else if (g.counter > 100 && g.counter < 200) {
			bgOptions := &ebiten.DrawImageOptions {}
			screen.DrawImage(g.bgImage, bgOptions)
			screen.DrawImage(g.waterImage, bgOptions)
			g.block = &blocks[1]
			g.blockX = (g.field.width/2) - 1
			g.blockY = g.field.height/2 - 2
			g.blockY -= (g.counter % 50)/1
			g.DrawBlock(screen, bgOptions)		
			handOptions := &ebiten.DrawImageOptions {}
			handOptions.GeoM.Translate(45, 20);
			if (g.counter % 20 > 10) {
				screen.DrawImage(g.handImage, handOptions)			
			}		
			text.Draw(screen, "To drop \nbubbbles up \ntap top side", smallFont, 10, 60, darkMagenta);
		} else if (g.counter > 200 && g.counter < 300) {
			bgOptions := &ebiten.DrawImageOptions {}
			screen.DrawImage(g.bgImage, bgOptions)
			screen.DrawImage(g.waterImage, bgOptions)
			g.block = &blocks[1]
			g.blockX = (g.field.width/2) - 1
			g.blockY = g.field.height/2 - 2
			g.blockX -= (g.counter % 50)/10
			g.DrawBlock(screen, bgOptions)
			handOptions := &ebiten.DrawImageOptions {}
			handOptions.GeoM.Translate(10, 90);
			if (g.counter % 20 > 10) {
				screen.DrawImage(g.handImage, handOptions)			
			}		
			text.Draw(screen, "To move \nbubbbles left \ntap left side", smallFont, 10, 140, magenta);
		} else if (g.counter > 300 && g.counter < 400) {
			bgOptions := &ebiten.DrawImageOptions {}
			screen.DrawImage(g.bgImage, bgOptions)
			screen.DrawImage(g.waterImage, bgOptions)
			g.block = &blocks[1]
			g.blockX = (g.field.width/2) - 1
			g.blockY = g.field.height/2	- 2			
			i := 1
			for i <= g.counter / 25 {
				i = i + 1
				g.Rotate()
			}			
			g.DrawBlock(screen, bgOptions)
			handOptions := &ebiten.DrawImageOptions {}
			handOptions.GeoM.Translate(45, 90);
			if (g.counter % 50 > 25) {
				screen.DrawImage(g.handImage, handOptions)			
			}	
			text.Draw(screen, "To rotate \nfigure clockwise \ntap it", smallFont, 10, 140, magenta);	
		}

		if g.splashPressed {
			vector.DrawFilledRect(screen, 10, 175, 60, 15, magenta, false)
			text.Draw(screen, "SKIP INTRO", smallFont, 13, 187, color.White)		
		} else {
			vector.DrawFilledRect(screen, 10, 175, 60, 15, darkMagenta, false)
			vector.DrawFilledRect(screen, 11, 176, 60, 15, magenta, false)
			text.Draw(screen, "SKIP INTRO", smallFont, 14, 186, color.White)		
		}

		return
	}

	if g.block == nil && !g.over {

		g.block = &blocks[rand.Intn(7)]
		g.blockX = 0
 		g.blockY = g.field.height - g.block.height - bottomHeight //to get space for counter
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
	bgOptions := &ebiten.DrawImageOptions {}
	smallWeedOptions := &ebiten.DrawImageOptions {}
	weedOptions := &ebiten.DrawImageOptions {}
	screen.Fill(colors[0])

	screen.DrawImage(g.bgImage, bgOptions)	

	weedOptions.GeoM.Translate(0, 60)
	
	sx, sy := 0, 0

	i := (g.counter / 10) % 15

	if (g.level > 1) {

	sx, sy := frameOX+i*frameWidth, frameOY
	screen.DrawImage(g.weedImage.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image), weedOptions)

	} 

	if (g.level > 2) {
		g.speed = 12
	} 	

	if (g.level > 3) {
		weedOptions.GeoM.Translate(40, 0)
		i = (i + 5) % 15
		sx, sy = frameOX+i*frameWidth, frameOY
		screen.DrawImage(g.weedImage.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image), weedOptions)
	}

	if (g.level > 4) {
		g.speed = 9
	}

	if (g.level > 5) {
		weedOptions.GeoM.Translate(-70, 0)
		i = (i + 5) % 15
		sx, sy = frameOX+i*frameWidth, frameOY
		screen.DrawImage(g.weedImage.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image), weedOptions)
	}

	if (g.level > 6) {
		g.speed = 6
	}

	screen.DrawImage(g.waterImage, bgOptions)

	// Draw block	
	if (!g.over) {

		g.DrawBlock(screen, options)

		if (g.BlockLanded()) {
			g.MergeBlock()
			g.RemoveLines()
		}

	}	

	smallWeedOptions.GeoM.Translate(15, 145)
	i = (i + 5) % 15
	sx, sy = frameOX+i*28, frameOY
	screen.DrawImage(g.smallWeedImage.SubImage(image.Rect(sx, sy, sx+28, sy+44)).(*ebiten.Image), smallWeedOptions)

	smallWeedOptions.GeoM.Translate(35, 0)
	i = (i + 5) % 15
	sx, sy = frameOX+i*28, frameOY
	screen.DrawImage(g.smallWeedImage.SubImage(image.Rect(sx, sy, sx+28, sy+44)).(*ebiten.Image), smallWeedOptions)

	screen.DrawImage(g.waterImage, bgOptions)

	vector.DrawFilledRect(screen, 0, 
		float32(g.cell_size * (g.field.height - bottomHeight)), 
		float32(g.field.width * g.cell_size), 
		float32(g.cell_size * bottomHeight), colors[10], false)
	text.Draw(screen, fmt.Sprintf("SCORE: %d  LEVEL: %d", g.lines, g.level), 
	smallFont, 9, 9 + g.cell_size * (g.field.height - bottomHeight), color.White)

	//screen.DrawImage(g.waterImage, bgOptions)

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
	g.freq = g.speed
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
			g.lines++
			fmt.Printf("g.lines = %d\n",  g.lines)
			if (g.lines < 2) {
				fmt.Println("A")
				g.level = 2
			} else if (g.lines < 4) {
				fmt.Println("B")
				g.level = 3
			} else if (g.lines < 8) {
				fmt.Println("C")
				g.level = 4
			} else if (g.lines < 16) {
				fmt.Println("D")
				g.level = 5
			} else {
				fmt.Println("E")
				g.level = 6
			} 
			g.RemoveLine(j)
		}
	}
}

func (g *Game) Update () error {

	//fmt.Println("Update")

	if (ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)) {
		x, y := ebiten.CursorPosition()
		if x > 20 && x < 80 && y > 100 && y < 120 && !g.splash {
			g.pressed = true
		} else if x > 10 && x < 80 && y > 175 && y < 190 && g.splash {
			g.splashPressed = true
		}
	} else if g.splashPressed {
		g.splashPressed = false
		g.splash = false
		g.block = nil
	} else if g.pressed {
		g.pressed = false
		g.over = false
		g.lines = 0
		g.level = 1
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
			if !g.over {
				if (g.Intersects(g.block, x, y)) {
					g.Rotate()  
				} else if (x < 50 && y > 50) {
					g.MoveLeft()
				} else if (x > 50 && y > 50) {
					g.MoveRight()				
				} else if (y < 50) {
					g.Drop()
				} 
			} else {
				fmt.Printf("pressed = %t over = %t\n", g.pressed, g.over)
				if (x > 20 && x < 80 && y > 100 && y < 120) {
					fmt.Println("A");
					g.pressed = true
				} else if g.pressed {
					fmt.Println("B");
					g.pressed = false
					g.over = false
					g.field = &Matrix{
						cells:make([]int, 10*20),
						width:10,
						height:20,			
					}		
				} else {
					fmt.Println("C");
					g.pressed = false
				}
			}

			if g.splash {
				fmt.Printf("pressed = %t over = %t\n", g.pressed, g.over)
				if (x > 10 && x < 80 && y > 175 && y < 190) {
					fmt.Println("A");
					g.splashPressed = true
				} else if g.splashPressed {
					fmt.Println("B");
					g.splashPressed = false
					g.splash = false
				} else {
					fmt.Println("C");
					g.splashPressed = false
				}
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
	return g.field.width * g.cell_size, g.field.height * g.cell_size
}
