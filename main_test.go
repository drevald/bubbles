package main

import (
	"testing"
	"fmt"
	"os"
	"image/png"
	_ "image/png"
	"bytes"
	"github.com/hajimehoshi/ebiten/v2"
)

func TestRotateMatrix (t *testing.T) {
	m := Matrix{cells:[]int{1,1,1,0,1,0}, width:3, height:2}
	fmt.Println(m)
	fmt.Println("Rotating right")
	m.RotateRight()
	fmt.Println(m)
	fmt.Println("Rotating left")
	m.RotateLeft()
	fmt.Println(m)
	fmt.Println("Testing")
}

func TestReadExternalImage (t *testing.T) {
	data, err := os.ReadFile("bubble.png")
	if err != nil {
		fmt.Println(err)
	}
	imageReader := bytes.NewReader(data)
	image, err := png.Decode(imageReader)
	if err != nil {
		fmt.Println(err)
	}
	bubbleImage := ebiten.NewImageFromImage(image)
	screen := ebiten.NewImage(100, 100)
	options := &ebiten.DrawImageOptions{}
	screen.DrawImage(bubbleImage, options)
}