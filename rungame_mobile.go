//go:build mobile

package ge

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/mobile"
)

func rungame(g ebiten.Game) error {
	mobile.SetGame(g)
	return nil
}
