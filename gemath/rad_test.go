package gemath

import (
	"math"
	"testing"
)

func TestRadPositive(t *testing.T) {
	tests := []struct {
		a    Rad
		want Rad
	}{
		{0, 0},
		{-1, 5.2831853072},
		{-math.Pi, math.Pi},
		{1, 1},
		{4 * math.Pi, 4 * math.Pi},
		{math.Pi, math.Pi},
	}

	for _, test := range tests {
		have := test.a.Positive()
		if !have.EqualApprox(test.want) {
			t.Fatalf("Positive(%f):\nhave: %.10f\nwant: %.10f", test.a, have, test.want)
		}
	}
}

func TestRadNormalized(t *testing.T) {
	tests := []struct {
		a    Rad
		want Rad
	}{
		{0, 0},
		{1, 1},
		{3 * math.Pi, math.Pi},
		{math.Pi, math.Pi},
		{-math.Pi, math.Pi},
		{-0.2, 2*math.Pi - 0.2},
	}

	for _, test := range tests {
		have := test.a.Normalized()
		if !have.EqualApprox(test.want) {
			t.Fatalf("Normalized(%f):\nhave: %f\nwant: %f", test.a, have, test.want)
		}
	}
}

func TestRadAngleDelta(t *testing.T) {
	tests := []struct {
		a    Rad
		b    Rad
		want Rad
	}{
		{0, 0, 0},
		{1, 1, 0},
		{-0.2, 0.2, 0.4},
		{0.4, 0.2, -0.2},
		{-0.5, -0.2, 0.3},
		{0.4, 0, -0.4},
		{math.Pi, 0, -math.Pi},
		{-math.Pi, 0, math.Pi},
		{2 * math.Pi, 0, 0},
		{4 * math.Pi, 0, 0},
		{6 * math.Pi, 0, 0},
		{3 * math.Pi, 0, -math.Pi},
		{0, 3 * math.Pi, math.Pi},
	}

	for _, test := range tests {
		have := test.a.AngleDelta(test.b)
		if !have.EqualApprox(test.want) {
			t.Fatalf("AngleDelta(%f, %f):\nhave: %f\nwant: %f",
				test.a, test.b, have, test.want)
		}
	}
}
