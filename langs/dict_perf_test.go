package langs

import "testing"

func BenchmarkDictGet(b *testing.B) {
	d := NewDictionary("en", 2)
	err := d.Load("", []byte(`
##first_part : example
##first_part.second_part : example
##first_part.second_part.third_part : example
##first_part.second_part.third_part.fourth_part : example
	`))
	if err != nil {
		b.Fatal(err)
	}

	parts := []string{
		"first_part",
		"second_part",
		"third_part",
		"fourth_part",
		"first_part.second_part.third_part.fourth_part",
	}

	b.Run("part1long", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			d.Get(parts[len(parts)-1])
		}
	})
	b.Run("part1long_const", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			d.Get("first_part.second_part.third_part.fourth_part")
		}
	})
	b.Run("part1", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			d.Get(parts[0])
		}
	})
	b.Run("part2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			d.Get(parts[0], parts[1])
		}
	})
	b.Run("part3", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			d.Get(parts[0], parts[1], parts[2])
		}
	})
	b.Run("part4", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			d.Get(parts[0], parts[1], parts[2], parts[3])
		}
	})
}
