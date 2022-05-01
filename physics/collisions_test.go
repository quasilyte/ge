package physics

import (
	"fmt"
	"testing"

	"github.com/quasilyte/ge/gemath"
)

func TestCircleCircleCollision(t *testing.T) {
	type testCircle struct {
		pos gemath.Vec
		r   float64
	}
	vec := func(x, y float64) gemath.Vec {
		return gemath.Vec{X: x, Y: y}
	}

	tests := []struct {
		a      testCircle
		b      testCircle
		want   []Collision
		config CollisionConfig
	}{
		{
			testCircle{pos: vec(0, 0), r: 1},
			testCircle{pos: vec(2, 2), r: 1},
			nil,
			CollisionConfig{},
		},

		{
			testCircle{pos: vec(0, 0), r: 1},
			testCircle{pos: vec(1, 1), r: 1},
			[]Collision{
				{LayerMask: 0b1},
			},
			CollisionConfig{},
		},
	}

	for i := range tests {
		test := tests[i]
		t.Run(fmt.Sprintf("test%d", i), func(t *testing.T) {
			var e CollisionEngine
			var body1 Body
			body1.InitCircle(nil, test.a.r)
			body1.Pos = test.a.pos
			var body2 Body
			body2.InitCircle(nil, test.b.r)
			body2.Pos = test.b.pos
			e.AddBody(&body1)
			e.AddBody(&body2)
			e.CalculateFrame()
			collisions := e.GetCollisions(&body1, CollisionConfig{})
			if len(test.want) != len(collisions) {
				t.Fatalf("%s vs %s: expected %d collisions, have %d", body1, body2, len(test.want), len(collisions))
			}
			for i := range collisions {
				have := collisions[i]
				want := test.want[i]
				if have.LayerMask != want.LayerMask {
					t.Fatalf("%s vs %s: collision layer mask mismatches:\nhavw: %b\nwant: %b",
						body1, body2, have.LayerMask, want.LayerMask)
				}
			}
		})
	}
}
