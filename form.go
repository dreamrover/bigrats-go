package main

import (
	"os/exec"
	"runtime"

	"github.com/visualfc/goqt/ui"
)

const (
	width  = 960
	height = 600
	cols   = 6
)

type Form struct {
	*ui.QWidget
	taskButton *ui.QPushButton
	label      [4]*ui.QLabel
	spinBox    *ui.QSpinBox
	checkBox   [2]*ui.QCheckBox
	comboBox   *ui.QComboBox
	table      [2]*Table
	tabWidget  *ui.QTabWidget
	timer      *ui.QTimer
}

func gui() {
	var hbox [2]*ui.QHBoxLayout
	var vbox [3]*ui.QVBoxLayout
	var err error

	columnLabels := [...]string{"Name", "Site", "Size", "Progress", "Speed", "ETA"}

	w := &Form{}
	w.QWidget = ui.NewWidget()

	for i := range hbox {
		hbox[i] = ui.NewHBoxLayout()
	}
	for i := range vbox {
		vbox[i] = ui.NewVBoxLayout()
	}
	for i := range w.label {
		w.label[i] = ui.NewLabel()
	}
	for i := range w.checkBox {
		w.checkBox[i] = ui.NewCheckBox()
	}
	for i := range w.table {
		w.table[i] = NewTable()
		w.table[i].SetEditTriggers(ui.QAbstractItemView_NoEditTriggers)
		w.table[i].SetSelectionBehavior(ui.QAbstractItemView_SelectRows)
		w.table[i].SetColumnCount(cols)
		w.table[i].SetHorizontalHeaderLabels(columnLabels[:])
		w.table[i].ResizeColumnsToContents()
		w.table[i].ResizeRowsToContents()
		w.table[i].SetColumnWidth(0, 200)
	}

	w.taskButton = ui.NewPushButton()
	w.taskButton.SetText("New Task")
	w.taskButton.OnClicked(func() {
		go runTask("")
	})

	w.label[0].SetText("Threads")

	w.spinBox = ui.NewSpinBox()
	w.spinBox.SetValue(threads)
	w.spinBox.SetRange(0, MaxThreads)
	w.spinBox.OnValueChangedWithInt32(func(n int32) {
		if n > threads {
			chTask <- nil
		}
		threads = n
	})

	w.checkBox[0] = ui.NewCheckBox()
	w.checkBox[0].SetText("Auto Merge")
	for _, s := range avidemux {
		merger, err = exec.LookPath(s)
		if err == nil {
			break
		}
	}
	if merger == "" {
		automerge = false
	}
	w.checkBox[0].SetChecked(automerge)
	w.checkBox[0].OnClickedEx(func(checked bool) {
		if checked {
			for _, s := range avidemux {
				merger, err = exec.LookPath(s)
				if err == nil {
					break
				}
			}
			if merger == "" {
				w.Msgbox(err.Error())
				automerge = false
				autodel = false
				w.checkBox[0].SetChecked(false)
				w.checkBox[1].SetChecked(false)
				w.checkBox[1].SetCheckable(false)
				w.comboBox.SetEnabled(false)
				return
			}
			automerge = true
			w.checkBox[1].SetCheckable(true)
			w.checkBox[1].SetChecked(autodel)
			autodel = true
			w.comboBox.SetEnabled(true)
		} else {
			automerge = false
			autodel = false
			w.checkBox[1].SetChecked(false)
			w.checkBox[1].SetCheckable(false)
			w.comboBox.SetEnabled(false)
		}
	})

	w.checkBox[1].SetText("Delete Segments")
	w.checkBox[1].SetChecked(autodel)
	if !w.checkBox[0].IsChecked() {
		w.checkBox[1].SetCheckable(false)
	}
	w.checkBox[1].OnClickedEx(func(checked bool) {
		autodel = checked
	})

	w.label[1].SetText("File Format:")

	w.comboBox = ui.NewComboBox()
	w.comboBox.AddItems([]string{"Original", ".mp4", ".flv", ".avi", ".mkv"})
	w.comboBox.SetCurrentIndex(cindex)
	w.comboBox.SetEditable(false)
	container = w.comboBox.ItemText(cindex)
	if container == "" {
		cindex = 0
		w.comboBox.SetCurrentIndex(cindex)
		container = w.comboBox.ItemText(cindex)
	}
	if !w.checkBox[0].IsChecked() {
		w.comboBox.SetEnabled(false)
	}
	w.comboBox.OnCurrentIndexChanged(func(s string) {
		container = s
		cindex = w.comboBox.CurrentIndex()
	})

	w.table[0].VerticalHeader().SetVisible(false)

	/*w.table[0].SetContextMenuPolicy(ui.Qt_CustomContextMenu)
	var action [2]*ui.QAction
	menu := ui.NewMenu()
	action[0] = ui.NewActionWithTextParent("pause", w)
	action[1] = menu.AddActionWithText("resume")
	w.table[0].OnCustomContextMenuRequested(func(point *ui.QPoint) {
		menu.AddAction(action[0])
		pos := w.table[0].MapToGlobal(point)
		menu.ExecWithPos(pos)
	})
	action[0].OnTriggered(func() {
		indexes := w.table[0].SelectedIndexes()
		for _, index := range indexes {
			w.table[0].rows[index.Row()].seg.status = PAUSE
		}
		w.Msgbox(strconv.Itoa(int(indexes[0].Row())))
	})*/

	hbox[0].AddWidget(w.taskButton)
	hbox[0].AddWidget(w.label[0])
	hbox[0].AddWidget(w.spinBox)
	hbox[0].AddWidget(w.checkBox[0])
	hbox[0].AddWidget(w.checkBox[1])
	hbox[0].AddWidget(w.label[1])
	hbox[0].AddWidget(w.comboBox)
	hbox[0].AddStretch()

	hbox[1].AddWidget(w.label[2])
	hbox[1].AddWidget(w.label[3])
	hbox[1].SetAlignment(ui.Qt_AlignRight)

	vbox[1].AddWidget(w.table[0])
	vbox[1].AddLayout(hbox[1])

	frame := ui.NewFrame()
	frame.SetLayout(vbox[1])

	vbox[2].AddWidget(w.table[1])

	widget := ui.NewWidget()
	widget.SetLayout(vbox[2])

	w.tabWidget = ui.NewTabWidget()
	w.tabWidget.AddTabWithWidgetString(frame, "Downloading")
	w.tabWidget.AddTabWithWidgetString(widget, "Finished")

	vbox[0].AddLayout(hbox[0])
	vbox[0].AddWidget(w.tabWidget)

	w.SetLayout(vbox[0])
	w.SetWindowTitle("bigrats for " + runtime.GOOS + " " + runtime.GOARCH)

	desktop := ui.NewDesktopWidget()
	dwidth := desktop.Width()
	dheight := desktop.Height()
	w.Geometry().SetWidth(width)
	w.Geometry().SetHeight(height)
	w.SetGeometryWithXYWidthHeight((dwidth-width)/2, (dheight-height)/2, width, height)

	w.timer = ui.NewTimer()
	w.timer.OnTimeout(func() {
		var speed rate
		for {
			select {
			case seg := <-chRow:
				if seg.status != DONE {
					w.table[0].Refresh(seg, false)
					speed += seg.speed
				} else {
					w.table[0].Refresh(seg, true)
					seg.eta = "-"
					w.table[1].Refresh(seg, false)
				}
			case url := <-chURL:
				w.GetDir(url)
			case msg := <-chMsg:
				w.Msgbox(msg)
			case merge := <-chMrg:
				w.label[2].SetText(merge)
			default:
				if speed > 0 {
					w.label[3].SetText(speed.String())
					speed = 0
				} else {
					w.label[3].SetText("")
				}
				return
			}
		}
	})
	w.timer.StartWithMsec(1000)

	w.Show()
}
