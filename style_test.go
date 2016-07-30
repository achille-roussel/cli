package cli

import "testing"

func TestStyles(t *testing.T) {
	tests := []struct {
		name  string
		style StyleSet
	}{
		{"normal", Normal},
		{"bold", Bold},
		{"dim", Dim},
		{"standout", Standout},
		{"underline", Underline},
		{"blink", Blink},
		{"reverse", Reverse},
		{"hidden", Hidden},

		{"black-fg", Black},
		{"red-fg", Red},
		{"green-fg", Green},
		{"yellow-fg", Yellow},
		{"blue-fg", Black},
		{"magenta-fg", Magenta},
		{"cyan-fg", Cyan},
		{"white-fg", White},
		{"default-fg", Default},

		{"black-bg", BlackBG},
		{"red-bg", RedBG},
		{"green-bg", GreenBG},
		{"yellow-bg", YellowBG},
		{"blue-bg", BlackBG},
		{"magenta-bg", MagentaBG},
		{"cyan-bg", CyanBG},
		{"white-bg", WhiteBG},
		{"default-bg", DefaultBG},
	}

	for _, test := range tests {
		t.Logf("%s: %s", test.name, test.style.S("Hello World!"))
	}
}

func TestColors(t *testing.T) {
	tests := []struct {
		name  string
		color StyleSet
	}{
		{"black", Black},
		{"red", Red},
		{"green", Green},
		{"yellow", Yellow},
		{"blue", Blue},
		{"magenta", Magenta},
		{"cyan", Cyan},
		{"white", White},
	}

	for _, test := range tests {
		t.Logf("%s: %s", test.name, test.color.S(test.name))
	}
}

func TestStripStyles(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{
			in:  "",
			out: "",
		},
		{
			in:  "Hello World!",
			out: "Hello World!",
		},
		{
			in:  "\033[1mHello World!\033[0m",
			out: "Hello World!",
		},
		{
			in:  "Hello \033[32mWorld!\033[0m",
			out: "Hello World!",
		},
	}

	for _, test := range tests {
		if s := string(StripStyles([]byte(test.in))); s != test.out {
			t.Errorf("%s: %#v != %#v", test.in, test.out, s)
		}
	}
}
