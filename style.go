package cli

import (
	"fmt"
	"io"
	"strconv"
	"unicode/utf8"
)

type StyleSet []int

func (s StyleSet) S(v string) string {
	return string(s.B([]byte(v)))
}

func (s StyleSet) B(b []byte) []byte {
	return applyStyleCodes(b, s...)
}

func (s StyleSet) Append(b []byte) []byte {
	return appendStyleCodes(b, s...)
}

func (s StyleSet) Bytes() []byte {
	return s.Append(make([]byte, 0, 4+2*len(s)))
}

func (s StyleSet) String() string {
	return string(s.Bytes())
}

func (s StyleSet) Format(f fmt.State, _ rune) {
	s.WriteTo(f)
}

func (s StyleSet) WriteTo(w io.Writer) (int, error) {
	return w.Write(s.Bytes())
}

func Style(styles ...StyleSet) StyleSet {
	set := make(StyleSet, 0, 2*len(styles))

	for _, s := range styles {
		set = append(set, s...)
	}

	return set
}

var (
	Normal    StyleSet = StyleSet{0}
	Bold      StyleSet = StyleSet{1}
	Dim       StyleSet = StyleSet{2}
	Standout  StyleSet = StyleSet{3}
	Underline StyleSet = StyleSet{4}
	Blink     StyleSet = StyleSet{5}
	Reverse   StyleSet = StyleSet{7}
	Hidden    StyleSet = StyleSet{8}

	Black   StyleSet = StyleSet{30}
	Red     StyleSet = StyleSet{31}
	Green   StyleSet = StyleSet{32}
	Yellow  StyleSet = StyleSet{33}
	Blue    StyleSet = StyleSet{34}
	Magenta StyleSet = StyleSet{35}
	Cyan    StyleSet = StyleSet{36}
	White   StyleSet = StyleSet{37}
	Default StyleSet = StyleSet{38}

	BlackBG   StyleSet = StyleSet{40}
	RedBG     StyleSet = StyleSet{41}
	GreenBG   StyleSet = StyleSet{42}
	YellowBG  StyleSet = StyleSet{43}
	BlueBG    StyleSet = StyleSet{44}
	MagentaBG StyleSet = StyleSet{45}
	CyanBG    StyleSet = StyleSet{46}
	WhiteBG   StyleSet = StyleSet{47}
	DefaultBG StyleSet = StyleSet{48}
)

func ForEachByteInString(s string, do func(byte)) {
	ForEachByte([]byte(s), do)
}

func ForEachByte(b []byte, do func(byte)) {
	for i, n := 0, len(b); i != n; {
		c := b[i]

		if c == '\033' {
			for i = i + 1; i != n; i++ {
				if b[i] == 'm' {
					i++
					break
				}
			}
			continue
		}

		do(c)
		i++
	}
}

func ForEachRuneInString(s string, do func(rune)) {
	ForEachRune([]byte(s), do)
}

func ForEachRune(b []byte, do func(rune)) {
	for i, n := 0, len(b); i != n; {
		c, z := utf8.DecodeRune(b[i:])

		if c == '\033' {
			for i = i + 1; i != n; i++ {
				if b[i] == 'm' {
					i++
					break
				}
			}
			continue
		}

		do(c)
		i += z
	}
}

func RuneCountInString(s string) int {
	return RuneCount([]byte(s))
}

func RuneCount(b []byte) (n int) {
	ForEachRune(b, func(_ rune) { n++ })
	return
}

func StripStylesInString(s string) string {
	return string(StripStyles([]byte(s)))
}

func StripStyles(b []byte) []byte {
	j := 0
	ForEachByte(b, func(c byte) {
		b[j] = c
		j++
	})
	return b[:j]
}

func appendStyleCodes(b []byte, codes ...int) []byte {
	b = append(b, '\033', '[')

	for i, c := range codes {
		if i != 0 {
			b = append(b, ';')
		}
		b = strconv.AppendInt(b, int64(c), 10)
	}

	return append(b, 'm')
}

func applyStyleCodes(b []byte, codes ...int) []byte {
	s := make([]byte, 0, len(b)+(4*len(codes)))
	s = appendStyleCodes(s, codes...)
	s = append(s, b...)
	s = appendStyleCodes(s, 0)
	return s
}
