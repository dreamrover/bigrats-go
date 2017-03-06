package main

import "github.com/visualfc/goqt/ui"

type TableRow struct {
	items *[cols]*ui.QTableWidgetItem
	bar   *ui.QProgressBar
	seg   *seginfo
}

type Table struct {
	*ui.QTableWidget
	rows map[[16]byte]*TableRow
}

func NewTable() *Table {
	return &Table{ui.NewTableWidget(), make(map[[16]byte]*TableRow)}
}

func NewTableRow() *TableRow {
	var items [cols]*ui.QTableWidgetItem
	for i, _ := range items {
		items[i] = ui.NewTableWidgetItem()
		if i > 0 {
			items[i].SetTextAlignment(0x82) // center right
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
	return n
}

func (t *Table) DeleteRow(id [16]byte) {
	i := t.rows[id].items[0].Row()
	for j, _ := range t.rows[id].items {
		t.RemoveCellWidget(i, int32(j))
	}
	t.rows[id].bar.Delete()
	for _, item := range t.rows[id].items {
		item.Delete()
	}
	delete(t.rows, id)
	t.RemoveRow(i)
}

func (t *Table) Refresh(seg *seginfo, delrow bool) {
	var row *TableRow

	row, ok := t.rows[seg.sid]
	if !ok {
		row = NewTableRow()
		t.rows[seg.sid] = row
		row.seg = seg
		t.AppendRow(row)
		row.items[0].SetText(seg.name)
		row.items[1].SetText(seg.site)
	}
	row.items[2].SetText(seg.size.String())
	row.items[3].SetText(seg.speed.String())
	row.items[4].SetText(seg.eta)

	if seg.status == ERROR {
		row.bar.SetMaximum(100)
		color := ui.NewColor()
		color.SetRed(180)
		p := row.bar.Palette()
		p.SetColorWithCrColor(ui.QPalette_Base, color)
		row.bar.SetPalette(p)
	} else if seg.size >= 0 {
		row.bar.SetMaximum(int32(seg.size))
	}
	if seg.size > 0 {
		row.bar.SetValue(int32(seg.down))
		if seg.status != DONE {
			color := ui.NewColor()
			color.SetGreen(180)
			p := row.bar.Palette()
			p.SetColorWithCrColor(ui.QPalette_Highlight, color)
			row.bar.SetPalette(p)
		}
	}

	t.ResizeColumnsToContents()

	if delrow {
		t.DeleteRow(row.seg.sid)
	}
}
