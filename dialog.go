package main

import (
	"strings"

	"github.com/visualfc/goqt/ui"
)

func (f *Form) Msgbox(msg string) {
	msgbox := ui.NewMessageBox()
	msgbox.SetParent(f)
	msgbox.SetText(msg)
	msgbox.Show()
}

func (f *Form) GetDir(scriptURL string) {
	var url, dir string
	var radio [2]*ui.QRadioButton
	var hbox [3]*ui.QHBoxLayout
	var label [2]*ui.QLabel

	for i := range hbox {
		hbox[i] = ui.NewHBoxLayout()
	}
	for i := range label {
		label[i] = ui.NewLabel()
		label[i].SetFixedWidth(50)
	}

	dialog := ui.NewDialog()
	dialog.SetParent(f)
	dialog.SetWindowTitle("New Task")

	filedialog := ui.NewFileDialog()
	filedialog.SetParent(dialog)
	filedialog.SetFileMode(ui.QFileDialog_DirectoryOnly)
	filedialog.SetDirectory(ui.QDirHome())
	if dirs.Len() > 0 {
		qdir := ui.NewDir()
		qdir.SetPath(dirs[0])
		if qdir.Exists() {
			filedialog.SetDirectory(qdir)
		}
	}

	radio[0] = ui.NewRadioButtonWithTextParent("Script URL", f)
	radio[1] = ui.NewRadioButtonWithTextParent("Xdown URL", f)
	if !xdown {
		radio[0].SetChecked(true)
	} else {
		radio[1].SetChecked(true)
	}
	radio[0].OnToggled(func(checked bool) {
		xdown = !checked
	})
	radio[1].OnToggled(func(checked bool) {
		xdown = checked
	})

	hbox[0].AddWidget(radio[0])
	hbox[0].AddWidget(radio[1])

	label[0].SetText("URL:")
	label[0].SetFixedWidth(50)
	line := ui.NewLineEdit()
	if scriptURL != "" {
		line.SetText(scriptURL)
		//line.SetReadOnly(true)
	}

	hbox[1].AddWidget(label[0])
	hbox[1].AddWidget(line)

	label[1].SetText("Folder:")
	comboBox := ui.NewComboBox()
	//comboBox.SetInsertPolicy(ui.QComboBox_InsertAtTop)
	comboBox.AddItem(filedialog.Directory().AbsolutePath())
	if dirs.Len() > 1 {
		comboBox.AddItems(dirs[1:dirs.Len()])
	}
	button := ui.NewPushButton()
	policy := ui.NewSizePolicy()
	policy.HorizontalStretch()
	button.SetSizePolicy(policy)
	button.SetText("Browse")
	button.OnClicked(func() {
		filedialog.Show()
	})
	filedialog.OnAccepted(func() {
		dir = filedialog.Directory().AbsolutePath()
		n := comboBox.FindText(dir)
		if n >= 0 {
			comboBox.SetCurrentIndex(n)
		} else {
			comboBox.InsertItems(0, []string{dir})
			comboBox.SetCurrentIndex(0)
		}
	})

	hbox[2].AddWidget(label[1])
	hbox[2].AddWidget(comboBox)
	hbox[2].AddWidget(button)

	checkBox := ui.NewCheckBox()
	checkBox.SetText("Create folder")
	checkBox.SetEnabled(false)

	buttonBox := ui.NewDialogButtonBox()
	buttonBox.SetStandardButtons(ui.QDialogButtonBox_Ok | ui.QDialogButtonBox_Cancel)
	buttonBox.SetCenterButtons(true)
	buttonBox.Button(ui.QDialogButtonBox_Ok).OnClicked(func() {
		if comboBox.CurrentText() != "" && line.Text() != "" {
			if !xdown {
				url = line.Text()
			} else {
				// http://www.flvcd.com/xdown.php?id=xxxxxxxx
				text := line.Text()
				i := strings.LastIndex(text, "?id=")
				url = "http://www.flvcd.com/diy/diy00" + text[i+4:] + ".htm"
			}
			dir = comboBox.CurrentText()
			dirs.Add(dir)
			chDir <- urldir{url, dir}
			dialog.Close()
			//dialog.Delete()
		}
	})
	buttonBox.Button(ui.QDialogButtonBox_Cancel).OnClicked(func() {
		url = ""
		dir = ""
		chDir <- urldir{url, dir}
		dialog.Delete()
	})

	vbox := ui.NewVBoxLayout()
	vbox.AddLayout(hbox[0])
	vbox.AddLayout(hbox[1])
	vbox.AddLayout(hbox[2])
	vbox.AddWidget(checkBox)
	vbox.AddWidget(buttonBox)
	dialog.SetLayout(vbox)
	dialog.SetFixedSizeWithWidthHeight(450, 200)
	dialog.SetWindowFlags(dialog.WindowFlags() | ui.Qt_Dialog)

	dialog.Show()
}
