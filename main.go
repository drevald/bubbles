package main

import (
	"github.com/drevald/bubbles/game"
	"github.com/hajimehoshi/ebiten/v2"
	"fmt"
	_ "image/png"	
)

func main() {		
	fmt.Println("Start Game")
	g := &game.Game{}
	g.Init()	
	ebiten.RunGame(g) 
}