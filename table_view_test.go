package cli

import (
	"bytes"
	"testing"
	"time"
)

func TestRenderTableViewFromStruct(t *testing.T) {
	type T struct {
		Question string    `table:"QUESTION:"`
		Answer   string    `table:"ANSWER:"`
		Date     time.Time `table:"DATE:"`
	}

	table := []T{
		T{
			Question: "How are you doing?",
			Answer:   "I'm fine thanks!",
			Date:     time.Now(),
		},
		T{
			Question: "What's your name?",
			Answer:   "Luke",
			Date:     time.Now(),
		},
		T{
			Question: "What time is it?",
			Answer:   "...",
			Date:     time.Now(),
		},
	}

	buffer := &bytes.Buffer{}
	buffer.Grow(4096)
	RenderTableView(buffer, StyledTableView(NewTableView(table), Bold))

	t.Logf("\n\n%s\n", buffer.String())
}

func TestRenderTableViewFromStructWithAnanymousFields(t *testing.T) {
	type A struct {
		Cool   bool `table:"COOL:"`
		Bedool bool `table:"BEDOOL:"`
	}

	type B struct {
		A
	}

	type T struct {
		B
		Question string    `table:"QUESTION:"`
		Answer   string    `table:"ANSWER:"`
		Date     time.Time `table:"DATE:"`
	}

	table := []T{
		T{
			B:        B{A{Cool: true}},
			Question: "How are you doing?",
			Answer:   "I'm fine thanks!",
			Date:     time.Now(),
		},
		T{
			B:        B{A{Cool: true}},
			Question: "What's your name?",
			Answer:   "Luke",
			Date:     time.Now(),
		},
		T{
			B:        B{A{Bedool: true}},
			Question: "What time is it?",
			Answer:   "...",
			Date:     time.Now(),
		},
	}

	buffer := &bytes.Buffer{}
	buffer.Grow(4096)
	RenderTableView(buffer, StyledTableView(NewTableView(table), Bold))

	t.Logf("\n\n%s\n", buffer.String())
}
