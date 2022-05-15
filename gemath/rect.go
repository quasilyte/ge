package gemath

type Rect struct {
	Min Vec
	Max Vec
}

func (r Rect) Width() float64  { return r.Max.X - r.Min.X }
func (r Rect) Height() float64 { return r.Max.Y - r.Min.Y }

func (r Rect) Center() Vec {
	return Vec{X: r.Max.X / 2, Y: r.Max.Y / 2}
}

func (r Rect) X1() float64 { return r.Min.X }
func (r Rect) Y1() float64 { return r.Min.Y }

func (r Rect) X2() float64 { return r.Max.X }
func (r Rect) Y2() float64 { return r.Max.Y }

func (r Rect) IsEmpty() bool {
	return r.Min.X >= r.Max.X || r.Min.Y >= r.Max.Y
}

func (r Rect) Contains(p Vec) bool {
	return r.Min.X <= p.X && p.X < r.Max.X &&
		r.Min.Y <= p.Y && p.Y < r.Max.Y

}

func (r Rect) Overlaps(other Rect) bool {
	return !r.IsEmpty() && !other.IsEmpty() &&
		r.Min.X < other.Max.X && other.Min.X < r.Max.X &&
		r.Min.Y < other.Max.Y && other.Min.Y < r.Max.Y
}
