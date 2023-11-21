package main

import (
	"github.com/drevald/bubbles/game"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	g := &game.Game{}
	g.Init()	
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.RunGame(g) 
}
