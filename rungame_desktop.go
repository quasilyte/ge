//go:build !mobile

package ge

import (
	"github.com/hajimehoshi/ebiten/v2"
)

func rungame(g ebiten.Game) error {
	return ebiten.RunGame(g)
}
