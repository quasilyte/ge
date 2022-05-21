package gemath

import (
	"sort"
)

// RandPicker performs a uniformly distributed random probing among the given objects with weights.
// Higher the weight, higher the chance of that object of being picked.
type RandPicker[T any] struct {
	r *Rand

	options   []randPickerOption[T]
	threshold float64
	sorted    bool
}

func NewRandPicker[T any](r *Rand) *RandPicker[T] {
	return &RandPicker[T]{r: r}
}

func (p *RandPicker[T]) Reset() {
	p.options = p.options[:0]
	p.threshold = 0
	p.sorted = false
}

func (p *RandPicker[T]) AddOption(value T, weight float64) {
	if weight == 0 {
		return // Zero probability in any case
	}
	p.threshold += weight
	p.options = append(p.options, randPickerOption[T]{
		value:     value,
		threshold: p.threshold,
	})
	p.sorted = false
}

func (p *RandPicker[T]) IsEmpty() bool {
	return len(p.options) != 0
}

func (p *RandPicker[T]) Pick() T {
	var result T
	if len(p.options) == 0 {
		return result // Zero value
	}
	if len(p.options) == 1 {
		return p.options[0].value
	}

	// In a normal use case the random picker is initialized and then used
	// without adding extra options, so this sorting will happen only once in that case.
	if !p.sorted {
		sort.Slice(p.options, func(i, j int) bool {
			return p.options[i].threshold < p.options[j].threshold
		})
		p.sorted = true
	}

	roll := p.r.FloatRange(0, p.threshold)
	i := sort.Search(len(p.options), func(i int) bool {
		return roll <= p.options[i].threshold
	})
	if i < len(p.options) && roll <= p.options[i].threshold {
		result = p.options[i].value
	} else {
		result = p.options[len(p.options)-1].value
	}
	return result
}

type randPickerOption[T any] struct {
	value     T
	threshold float64
}
