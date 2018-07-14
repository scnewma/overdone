package optional

import "testing"

func TestOrElse(t *testing.T) {
	cases := []struct {
		given    interface{}
		def      interface{}
		expected interface{}
	}{
		{"blah", "hello", "blah"},
		{3, 5, 3},
		{nil, "default", "default"},
	}

	for _, c := range cases {
		optional := Of(c.given)
		got := optional.OrElse(c.def)
		if got != c.expected {
			t.Errorf("OptionalOf(%q).OrElse(%q) == %q, want %q", c.given, c.def, got, c.expected)
		}
	}
}

func TestIsPresent(t *testing.T) {
	cases := []struct {
		given    interface{}
		expected bool
	}{
		{"blah", true},
		{3, true},
		{nil, false},
	}

	for _, c := range cases {
		opt := Of(c.given)
		got := opt.IsPresent()
		if got != c.expected {
			t.Errorf("given %q, IsPresent() == %t, want %t", c.given, got, c.expected)
		}
	}
}
