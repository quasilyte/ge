package main

import (
	"fmt"
	"math"
	"strings"

	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/input"
	"github.com/quasilyte/gmath"
)

var damageMaskImages = []resource.ImageID{
	ImageDamageMask1,
	ImageDamageMask2,
	ImageDamageMask3,
	ImageDamageMask4,
}

func formattedActionString(h *input.Handler, a input.Action) string {
	return formatKeys(h.ActionKeyNames(a, h.DefaultInputMask()))
}

func formatKeys(keys []string) string {
	result := make([]string, len(keys))
	for i, k := range keys {
		switch k {
		case "gamepad_a", "gamepad_start":
			result[i] = "(" + strings.TrimPrefix(k, "gamepad_") + ")"
		case "enter", "escape", "down", "up", "left", "right":
			result[i] = "[" + k + "]"
		default:
			panic(fmt.Sprintf("unexpected key: %q", k))
		}
	}
	return strings.Join(result, " or ")
}

func getLevelData(scene *ge.Scene, level int) []byte {
	return scene.LoadRaw(resource.RawID(int(RawLevel0JSON) + level)).Data
}

func createExplosions(scene *ge.Scene, pos gmath.Vec, min, max int) {
	numExplosions := scene.Rand().IntRange(min, max)
	for i := 0; i < numExplosions; i++ {
		offset := scene.Rand().Offset(-24, 24)
		e := newExplosion(pos.Add(offset))
		if scene.Rand().Chance(0.4) {
			e.Image = ImageExplosion3
		} else {
			e.Image = ImageExplosion1
		}
		scene.AddObject(e)
	}
}

func fireTargetPos(pos gmath.Vec, facing gmath.Rad, maxRange float64) gmath.Vec {
	switch facing {
	case facingRight:
		pos.X += maxRange
	case facingDown:
		pos.Y += maxRange
	case facingLeft:
		pos.X -= maxRange
	case facingUp:
		pos.Y -= maxRange
	}
	return pos
}

func facingTowards(pos, targetPos gmath.Vec) (gmath.Rad, bool) {
	var targetFacing gmath.Rad
	if pos.Y == targetPos.Y {
		if pos.X < targetPos.X {
			targetFacing = facingRight
		} else {
			targetFacing = facingLeft
		}
	} else if pos.X == targetPos.X {
		if pos.Y < targetPos.Y {
			targetFacing = facingDown
		} else {
			targetFacing = facingUp
		}
	} else {
		return 0, false
	}
	return targetFacing, true
}

func roundedPos(pos gmath.Vec) gmath.Vec {
	x := math.Floor(pos.X/64) * 64
	y := math.Floor(pos.Y/64) * 64
	return gmath.Vec{X: x + 32, Y: y + 32}
}

func fixedPos(pos gmath.Vec) gmath.Vec {
	if pos.X < 0 {
		pos.X = 32
	} else if pos.X > 1920 {
		pos.X = 1920 - 32
	}
	if pos.Y < 0 {
		pos.Y = 32
	} else if pos.Y > 896 {
		pos.Y = 896 - 32
	}
	return pos
}

func spriteHue(alliance int) gmath.Rad {
	if alliance == 1 {
		return 1.106539
	}
	return 0
}

const (
	buildingLayerMask = 1 << 2
	wallLayerMask     = 1 << 3
	mineLayer         = 1 << 4
	blockingLayerMask = buildingLayerMask | wallLayerMask
)

func unitLayerMask(alliance int) uint16 {
	return 1 << uint16(alliance)
}

func projectileLayerMask(alliance int) uint16 {
	return 0b11 ^ unitLayerMask(alliance)
}

func cellDelta(facing gmath.Rad) gmath.Vec {
	switch facing {
	case facingRight:
		return gmath.Vec{X: 64}
	case facingDown:
		return gmath.Vec{Y: 64}
	case facingLeft:
		return gmath.Vec{X: -64}
	default:
		return gmath.Vec{Y: -64}
	}
}
