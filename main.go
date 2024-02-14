package main

import (
	"github.com/drevald/bubbles/game"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	// splash := &game.Splash{}
	// ebiten.RunGame(splash)
	g := &game.Game{}
	g.Init()	
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.RunGame(g) 
}
