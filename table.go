package cli

import "fmt"

type Table struct {
	cols []string
	rows [][]string
}

func NewTable(columns ...string) *Table {
	return &Table{
		cols: columns,
	}
}

func NewTableFromView(view TableView) *Table {
	cols, rows := view.Size()

	t := &Table{
		cols: make([]string, cols),
		rows: make([][]string, rows),
	}

	for i := 0; i != cols; i++ {
		t.cols[i] = view.Column(i)
	}

	for j := 0; j != rows; j++ {
		row := make([]string, cols)

		for i := 0; i != cols; i++ {
			row[i] = view.Cell(i, j)
		}

		t.rows[j] = row
	}

	return t
}

func (t *Table) Append(row ...string) *Table {
	if len(t.cols) != len(row) {
		panic(fmt.Sprintf("cli.(*Table).Append: invalid row length (expected %d but found %d)", len(t.cols), len(row)))
	}

	t.rows = append(t.rows, row)
	return t
}

func (t *Table) Column(col int) string {
	return t.cols[col]
}

func (t *Table) Cell(col int, row int) string {
	return t.rows[row][col]
}

func (t *Table) Size() (cols int, rows int) {
	cols, rows = len(t.cols), len(t.rows)
	return
}
