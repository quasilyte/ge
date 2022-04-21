package gemath

type Rect struct {
	Min Vec
	Max Vec
}

func (r Rect) Width() float64  { return r.Max.X - r.Min.X }
func (r Rect) Height() float64 { return r.Max.Y - r.Min.Y }

func (r Rect) X1() float64 { return r.Min.X }
func (r Rect) Y1() float64 { return r.Min.Y }

func (r Rect) X2() float64 { return r.Max.X }
func (r Rect) Y2() float64 { return r.Max.Y }

func (r Rect) IsEmpty() bool {
	return r.Min.X >= r.Max.X || r.Min.Y >= r.Max.Y
}
