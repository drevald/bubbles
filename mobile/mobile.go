package mobile

import (
	"github.com/hajimehoshi/ebiten/v2/mobile"
	"github.com/drevald/bubbles/game"
)

func init() {
	game := &game.Game{}
	game.Init()	
	mobile.SetGame(game)
}

// Dummy is a dummy exported function.
//
// gomobile doesn't compile a package that doesn't include any exported function.
// Dummy forces gomobile to compile this package.
func Dummy() {}
