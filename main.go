package main

import (
	"github.com/drevald/bubbles"
	"github.com/hajimehoshi/ebiten/v2"
	"fmt"
	"os"
	"bytes"
	"image/png"
	_ "image/png"		
)

func main() {
		
	fmt.Println("Start Game")

	data, err := os.ReadFile("bubble.png")
	if err != nil {
		fmt.Println(err)
	}
	imageReader := bytes.NewReader(data)
	image, err := png.Decode(imageReader)
	if err != nil {
		fmt.Println(err)
	}
	
	ebiten.RunGame(&Game{
		field: &Matrix{
			cells:make([]int, 10*20),
			width:10,
			height:20,			
		}, 
		cell_size: 10,	
		freq:10,
		over:false,
		bubbleImage:ebiten.NewImageFromImage(image),
	}) 

}