package cli

import (
	"fmt"
	"io"
	"reflect"
)

type TableView interface {
	Column(col int) string

	Cell(col int, row int) string

	Size() (cols int, rows int)
}

func RenderTableView(w io.Writer, t TableView) (err error) {
	cols, rows := t.Size()
	return renderTableView(w, t, computeTableColumnWidths(t, cols, rows), cols, rows)
}

func renderTableView(w io.Writer, t TableView, widths []int, cols int, rows int) (err error) {
	if err = renderTableViewColumns(w, t, widths, cols, rows); err == nil {
		err = renderTableViewRows(w, t, widths, cols, rows)
	}
	return
}

func RenderTableViewColumns(w io.Writer, t TableView) (err error) {
	cols, rows := t.Size()
	return renderTableViewColumns(w, t, computeTableColumnWidths(t, cols, rows), cols, rows)
}

func renderTableViewColumns(w io.Writer, t TableView, widths []int, cols int, rows int) (err error) {
	for i := 0; i != cols; i++ {
		if i != 0 {
			if _, err = io.WriteString(w, " "); err != nil {
				return
			}
		}

		if err = RenderColumn(w, t.Column(i), widths[i]); err != nil {
			return
		}
	}

	_, err = io.WriteString(w, "\n")
	return
}

func RenderTableViewRows(w io.Writer, t TableView) (err error) {
	cols, rows := t.Size()
	return renderTableViewRows(w, t, computeTableColumnWidths(t, cols, rows), cols, rows)
}

func renderTableViewRows(w io.Writer, t TableView, widths []int, cols int, rows int) (err error) {
	for j := 0; j != rows; j++ {
		for i := 0; i != cols; i++ {
			if i != 0 {
				if _, err = io.WriteString(w, " "); err != nil {
					return
				}
			}

			if err = RenderCell(w, t.Cell(i, j), widths[i], column(t.Column(i)).alignment()); err != nil {
				return
			}
		}

		if _, err = io.WriteString(w, "\n"); err != nil {
			return
		}
	}
	return
}

func computeTableColumnWidths(t TableView, cols int, rows int) (widths []int) {
	widths = make([]int, cols)

	for i := 0; i != cols; i++ {
		widths[i] = RuneCountInString(column(t.Column(i)).string())
	}

	for j := 0; j != rows; j++ {
		for i := 0; i != cols; i++ {
			if w := RuneCountInString(t.Cell(i, j)); w > widths[i] {
				widths[i] = w
			}
		}
	}

	return
}

func NewTableView(v interface{}) TableView {
	if view, ok := v.(TableView); ok {
		return view
	}
	return newTableView(reflect.ValueOf(v))
}

func newTableView(v reflect.Value) TableView {
	t := v.Type()

	switch t.Kind() {
	case reflect.Array, reflect.Slice:
		switch t.Elem().Kind() {
		case reflect.Struct:
			return makeStructSliceTableView(v, t)
		}
	}

	panic(fmt.Sprintf("cli.NewTableView: unsupported value, expectes slice of struct or cli.TableView but got %T", v.Interface()))
}

type structSliceTableView struct {
	v reflect.Value
	t reflect.Type
	c []string
	f []func(reflect.Value) reflect.Value
}

func makeStructSliceTableView(v reflect.Value, t reflect.Type) structSliceTableView {
	e := t.Elem()
	s := structSliceTableView{
		v: v,
		t: t,
		c: make([]string, 0, e.NumField()),
	}

	forEachColumn(e, func(col string, get func(reflect.Value) reflect.Value) {
		s.c = append(s.c, col)
		s.f = append(s.f, get)
	})

	return s
}

func (s structSliceTableView) Column(col int) string {
	return s.c[col]
}

func (s structSliceTableView) Cell(col int, row int) (cell string) {
	cell = fmt.Sprint(s.f[col](s.v.Index(row)).Interface())
	return
}

func (s structSliceTableView) Size() (cols int, rows int) {
	cols, rows = len(s.c), s.v.Len()
	return
}

func forEachColumn(t reflect.Type, do func(string, func(reflect.Value) reflect.Value)) {
	for i, n := 0, t.NumField(); i != n; i++ {
		f := t.Field(i)
		g := func(v reflect.Value) reflect.Value {
			return v.FieldByIndex(f.Index)
		}

		if f.Anonymous {
			// We want to flatten out the anonymous fields so they appear in the
			// table as top-level columns.
			// We call recursively and decorate the callback to do the recusive
			// field lookup among multiple levels if necessary.
			forEachColumn(f.Type, func(col string, get func(reflect.Value) reflect.Value) {
				do(col, func(v reflect.Value) reflect.Value { return get(g(v)) })
			})
			continue
		}

		switch name := f.Tag.Get("table"); name {
		case "-":
		case "":
			do(f.Name, g)
		default:
			do(name, g)
		}
	}
}

func StyledTableView(t TableView, style StyleSet) TableView {
	return styledTableView{t, style}
}

type styledTableView struct {
	TableView
	style StyleSet
}

func (t styledTableView) Column(col int) string {
	s := t.TableView.Column(col)
	l := ""
	r := ""

	if n := len(s); n != 0 && s[0] == ':' {
		l, s = s[:1], s[1:]
	}

	if n := len(s); n != 0 && s[n-1] == ':' {
		s, r = s[:n-1], s[n-1:]
	}

	return l + t.style.S(s) + r
}
