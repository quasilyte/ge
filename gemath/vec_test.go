package gemath

import (
	"testing"
)

func TestVecAPI(t *testing.T) {
	assertTrue := func(v bool) {
		t.Helper()
		if !v {
			t.Fatal("assertion failed")
		}
	}

	// Make sure that zero values can be used as literals.
	// Also make sure that we can use *Result methods on rvalue.

	assertTrue(Vec{}.EqualApprox(Vec{}))
	assertTrue(Vec{}.IsZero())
	assertTrue(Vec{}.Len() == 0)
	assertTrue(Vec{X: 1}.Neg() == Vec{X: -1})

	// A special case.
	assertTrue(Vec{}.Normalized() == Vec{})
}

func TestVecNormalized(t *testing.T) {
	tests := []struct {
		v    Vec
		want Vec
	}{
		{Vec{1, 0}, Vec{1, 0}},
		{Vec{-1, 0}, Vec{-1, 0}},
		{Vec{0, 1}, Vec{0, 1}},
		{Vec{0, -1}, Vec{0, -1}},
		{Vec{3, 0}, Vec{1, 0}},
		{Vec{0, 3}, Vec{0, 1}},
		{Vec{1, 1}, Vec{0.70710678118654, 0.70710678118654}},
		{Vec{10, 13}, Vec{0.6097107608, 0.7926239891}},
	}

	for _, test := range tests {
		have := test.v.Normalized()
		if !have.EqualApprox(test.want) {
			t.Fatalf("Normalized(%s):\nhave: %v\nwant: %v", test.v, have, test.want)
		}
		have2 := test.v.Divf(test.v.Len())
		if !have.EqualApprox(have2) {
			t.Fatalf("div+len of %s:\nhave: %v\nwant: %v", test.v, have, test.want)
		}
		if !have.IsNormalized() {
			t.Fatalf("IsNormalized(Normalized(%s)) returned false", test.v)
		}
	}
}

func TestVecLen(t *testing.T) {
	tests := []struct {
		v    Vec
		want float64
	}{
		{Vec{}, 0},
		{Vec{1, 0}, 1},
		{Vec{0, 1}, 1},
		{Vec{1, 1}, 1.414213562373},
		{Vec{2, 1}, 2.236067977499},
		{Vec{-1, 0}, 1},
		{Vec{0, -1}, 1},
	}

	for _, test := range tests {
		have := test.v.Len()
		if !EqualApprox(have, test.want) {
			t.Fatalf("Len(%s):\nhave: %v\nwant: %v", test.v, have, test.want)
		}
	}
}

func TestVecEqualApprox(t *testing.T) {
	tests := []struct {
		a    Vec
		b    Vec
		want bool
	}{
		{Vec{}, Vec{}, true},
		{Vec{}, Vec{1, 1}, false},
		{Vec{1, 1}, Vec{1, 1}, true},
		{Vec{0.5, 0.1}, Vec{-1, -0.3}, false},
		{Vec{0.01, 0.01}, Vec{}, false},
		{Vec{1, 1}, Vec{1 + Epsilon/2, 1 - Epsilon/2}, true},
		{Vec{0, 0 + Epsilon}, Vec{}, true},
		{Vec{0, 0 - Epsilon}, Vec{}, true},
		{Vec{0.000000001, 0}, Vec{}, true},
		{Vec{0.0000000001, 0}, Vec{}, true},
	}

	for _, test := range tests {
		have := test.a.EqualApprox(test.b)
		if have != test.want {
			t.Fatalf("EqualApprox(%s, %s):\nhave: %v\nwant: %v", test.a, test.b, have, test.want)
		}
		have2 := test.b.EqualApprox(test.a)
		if have2 != test.want {
			t.Fatalf("EqualApprox(%s, %s):\nhave: %v\nwant: %v", test.b, test.a, have2, test.want)
		}
	}
}

func TestVecAdd(t *testing.T) {
	tests := []struct {
		a    Vec
		b    Vec
		want Vec
	}{
		{Vec{}, Vec{}, Vec{}},
		{Vec{1, 1}, Vec{}, Vec{1, 1}},
		{Vec{}, Vec{1, 1}, Vec{1, 1}},
		{Vec{1, 1}, Vec{1, 1}, Vec{2, 2}},
		{Vec{0.5, 0.1}, Vec{-1, -0.3}, Vec{-0.5, -0.2}},
	}

	for _, test := range tests {
		have := test.a.Add(test.b)
		if !have.EqualApprox(test.want) {
			t.Fatalf("Add(%s, %s):\nhave: %s\nwant: %s", test.a, test.b, have, test.want)
		}
		have2 := test.b.Add(test.a)
		if !have2.EqualApprox(test.want) {
			t.Fatalf("Add(%s, %s):\nhave: %s\nwant: %s", test.b, test.a, have2, test.want)
		}
	}
}

func TestVecNeg(t *testing.T) {
	tests := []struct {
		arg  Vec
		want Vec
	}{
		{Vec{0, 0}, Vec{0, 0}},
		{Vec{1, 1}, Vec{-1, -1}},
		{Vec{-1, 2}, Vec{1, -2}},
		{Vec{1.5, 0.5}, Vec{-1.5, -0.5}},
		{Vec{-1.5, -0.5}, Vec{1.5, 0.5}},
	}

	for _, test := range tests {
		have := test.arg.Neg()
		if !have.EqualApprox(test.want) {
			t.Fatalf("Neg(%s):\nhave: %s\nwant: %s", test.arg, have, test.want)
		}
	}
}
