package mobile

import (
	"github.com/hajimehoshi/ebiten/v2/mobile"
	"github.com/drevald/bubbles/game"
	//"github.com/hajimehoshi/go-inovation/ino"
)

func init() {
	game := &game.Game{}
	game.Init()
	mobile.SetGame(game)
	//inogame, err := ino.NewGame()
	//if err != nil {
	 	//panic(err)
	 //}
	//mobile.SetGame(inogame)
}

// Dummy is a dummy exported function.
//
// gomobile doesn't compile a package that doesn't include any exported function.
// Dummy forces gomobile to compile this package.
func Dummy() {}
