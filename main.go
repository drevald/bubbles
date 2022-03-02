package main

import (
	"github.com/drevald/bubbles/game"
	"github.com/hajimehoshi/ebiten/v2"
	"fmt"
)

func main() {
	fmt.Println("Start Game")
	ebiten.RunGame(game.NewGame())
}