package gemath

// Slider is a value that can be increased and decreased with
// a custom overflow/underflow behavior.
//
// Min/Max fields control the range of the accepted values.
//
// It's a useful foundation for more high-level concepts
// like progress bars and gauges, paginators, option button selector, etc.
type Slider struct {
	min   int
	max   int
	value int

	// Clamp makes the slider use clamping overflow/underflow strategy
	// instead of the default wrapping around strategy.
	Clamp bool
}

// SetBounds sets the slider values range.
// It also sets the current value to min.
// Use TrySetValue if you need to override that.
//
// If max<min, this method panics.
func (s *Slider) SetBounds(min, max int) {
	if s.max < s.min {
		panic("min bound value should be less than max")
	}
	s.min = min
	s.max = max
	s.value = min
}

// TrySetValue assigns the v value to the slider if it fits its range.
// Returns true whether the slider value was assigned.
func (s *Slider) TrySetValue(v int) bool {
	if v >= s.min && v <= s.max {
		s.value = v
		return true
	}
	return false
}

// Value returns the current slider value.
// The returned value is guarandeed to be in [min, max] range.
func (s *Slider) Value() int { return s.value }

// Len returns the range of the values.
// Basically, it returns the max-min result.
func (s *Slider) Len() int { return s.max - s.min }

// Sub subtracts v from the current slider value.
// The overflow/underflow behavior depends on the slider settings.
func (s *Slider) Sub(v int) {
	s.Add(0 - v)
}

// Add adds v to the current slider value.
// The overflow/underflow behavior depends on the slider settings.
func (s *Slider) Add(v int) {
	switch v {
	case 1:
		s.Inc()
	case -1:
		s.Dec()
	default:
		if s.Clamp {
			s.value = Clamp(s.value+v, s.min, s.max)
		} else {
			value := s.value + v
			l := s.max - s.min + 1
			if value < s.min {
				value += l * ((s.min-value)/l + 1)
			}
			s.value = s.min + (value-s.min)%l
		}
	}
}

// Dec subtracts 1 from the slider value.
// Semantically identical to Sub(1), but more efficient.
func (s *Slider) Dec() {
	s.addBit(-1)
}

// Inc adds 1 to the slider value.
// Semantically identical to Add(1), but more efficient.
func (s *Slider) Inc() {
	s.addBit(1)
}

func (s *Slider) addBit(v int) {
	value := s.value + v
	if value < s.min {
		if s.Clamp {
			value = s.min
		} else {
			value = s.max
		}
	} else if value > s.max {
		if s.Clamp {
			value = s.max
		} else {
			value = s.min
		}
	}
	s.value = value
}
