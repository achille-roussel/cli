package cli

import (
	"bytes"
	"testing"
)

func TestRenderTable(t *testing.T) {
	tests := []struct {
		name  string
		table *Table
	}{
		{
			name:  "Empty Table",
			table: NewTable(),
		},
		{
			name:  "Only Columns",
			table: NewTable("Column A", "Column B", "Column C"),
		},
		{
			name:  "Default Alignment",
			table: NewTable("Column A", "Column B", "Column C").Append("1", "2", "3"),
		},
		{
			name:  "Left Alignment",
			table: NewTable("Column A:", "Column B:", "Column C:").Append("1", "2", "3"),
		},
		{
			name:  "Right Alignment",
			table: NewTable(":Column A", ":Column B", ":Column C").Append("1", "2", "3"),
		},
		{
			name:  "Centered",
			table: NewTable(":Column A:", ":Column B:", ":Column C:").Append("1", "2", "3"),
		},
	}

	buffer := &bytes.Buffer{}
	buffer.Grow(4096)

	for _, test := range tests {
		buffer.Reset()

		if err := RenderTableView(buffer, test.table); err != nil {
			t.Error(err)
		} else {
			t.Logf("%s\n\n%s\n", test.name, buffer.String())
		}
	}
}
