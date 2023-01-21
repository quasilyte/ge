package langs

import "testing"

func TestDictGet(t *testing.T) {
	d := NewDictionary("en")
	err := d.Load("", []byte(`
##first_part : example1
##first_part.second_part : example2
##first_part.second_part.third_part : example3
##first_part.second_part.third_part.fourth_part : example4
	`))
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		keys   []string
		result string
	}{
		{nil, "{{}}"},
		{[]string{}, "{{}}"},
		{[]string{""}, "{{}}"},
		{[]string{"wrong_key"}, "{{wrong_key}}"},
		{[]string{"first_part.wrong_key"}, "{{first_part.wrong_key}}"},

		{[]string{"first_part"}, "example1"},
		{[]string{"first_part", "second_part"}, "example2"},
		{[]string{"first_part", "second_part", "third_part"}, "example3"},
		{[]string{"first_part", "second_part", "third_part", "fourth_part"}, "example4"},

		{[]string{"first_part.second_part"}, "example2"},
		{[]string{"first_part.second_part.third_part"}, "example3"},
		{[]string{"first_part.second_part.third_part.fourth_part"}, "example4"},
	}

	for _, test := range tests {
		result := d.Get(test.keys...)
		if result != test.result {
			t.Fatalf("Get(%q):\nhave: %q\nwant: %q", test.keys, result, test.result)
		}
	}

}
