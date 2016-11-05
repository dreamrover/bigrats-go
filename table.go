package main

import "github.com/visualfc/goqt/ui"

type TableRow struct {
	items *[cols]*ui.QTableWidgetItem
	bar   *ui.QProgressBar
	seg   *seginfo
}

type Table struct {
	*ui.QTableWidget
	rows []*TableRow
}

func NewTable() *Table {
	return &Table{ui.NewTableWidget(), nil}
}

func NewTableRow() *TableRow {
	var items [cols]*ui.QTableWidgetItem
	for i, _ := range items {
		items[i] = ui.NewTableWidgetItem()
		if i > 0 {
			items[i].SetTextAlignment(0x82) // align center right
		}
	}
	bar := ui.NewProgressBar()
	return &TableRow{&items, bar, nil}
}

func (t *Table) AppendRow(row *TableRow) int32 {
	n := t.RowCount()
	t.InsertRow(n)
	t.SetItem(n, 0, row.items[0])
	t.SetItem(n, 1, row.items[1])
	t.SetItem(n, 2, row.items[2])
	//row.bar.SetFormat("%p%, %v/%m")
	t.SetCellWidget(n, 3, row.bar)
	t.SetItem(n, 4, row.items[3])
	t.SetItem(n, 5, row.items[4])
	t.rows = append(t.rows[:n], row)
	return n
}

func (t *Table) DeleteRow(i int32) {
	if int(i) >= len(t.rows) {
		return
	}
	var j int32
	c := t.ColumnCount()
	for j = 0; j < c; j++ {
		t.RemoveCellWidget(i, j)
	}
	t.rows[i].bar.Delete()
	//t.RemoveRow(i)
	//copy(t.rows[i:], t.rows[i+1:len(t.rows)])
	//t.RemoveCellWidget(i, 3)
	t.SetRowHidden(i, true)
}

func (t *Table) Refresh(ss rowinfo) {
	var row, r *TableRow
	var i int

	if len(t.rows) > 0 {
		for i, r = range t.rows {
			if r.seg.sid == ss.seg.sid {
				row = r
				break
			}
		}
	}
	if row == nil {
		row = NewTableRow()
		row.seg = ss.seg
		i = int(t.AppendRow(row))
	}
	row.items[0].SetText(ss.seg.name)
	row.items[1].SetText(ss.seg.site)
	row.items[2].SetText(ss.size.String())
	row.items[3].SetText(ss.speed.String())
	row.items[4].SetText(ss.eta)

	if ss.size >= 0 {
		row.bar.SetMaximum(int32(ss.size))
	}
	row.bar.SetValue(int32(ss.down))
	/*color := ui.NewColor()
	color.SetGreen(0)
	p := bar.Palette()
	p.SetColorWithCrColor(ui.QPalette_Base, color)
	bar.SetPalette(p)*/
	//bar.SetStyleSheet()
	//bar.SetBackgroundRole(color.Green())

	t.ResizeColumnsToContents()

	if ss.status == "Delete" {
		t.DeleteRow(int32(i))
	}
}
