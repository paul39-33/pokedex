package main

import "testing"

func TestCleanInput(t *testing.T){
	cases := []struct {
		input string
		expected []string
	}{
		{
			input:	"  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:	"  ooooo 		yeahhh ",
			expected: []string{"ooooo", "yeahhh"},
		},
		{
			input: "Charmander Bulbasaur PIKACHU",
			expected: []string{"charmander", "bulbasaur", "pikachu"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected){
			t.Errorf("actual length: %v doesnt matched expected length: %v", len(actual), len(c.expected))
		}
		
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]

			if word != expectedWord{
				t.Errorf("actual word: %v doesn't match expected word: %v", word, expectedWord)
			}
		}
	}
}