package bubblemobile

import (
	"github.com/drevald/bubbles/game"
	"github.com/hajimehoshi/ebiten/v2/mobile"
)

func init() {
  mobile.SetGame(&game.Game{})
}

func Dummy() {

}