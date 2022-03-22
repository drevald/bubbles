package bubblemobile

import (
	"github.com/drevald/bubbles/game"
	"github.com/hajimehoshi/ebiten/v2/mobile"
	_ "image/png"
)

func init() {
	g := &game.Game{}
	g.Init()
	mobile.SetGame(g)
}

func Dummy() {

}