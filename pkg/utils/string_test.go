package utils

import "testing"

func TestIsBlank(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"", true},
		{"  ", true},
		{"a", false},
	}

	for _, c := range cases {
		got := IsBlank(c.in)
		if got != c.want {
			t.Errorf("IsBlank(%q) == %t, want %t", c.in, got, c.want)
		}
	}
}
