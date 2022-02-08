/*
Copyright © 2021 - 2022 Gianni Doria (gianni.doria@gmail.com)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cfg

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

type Item struct {

	// name: text
	Name    string            `yaml:"name"`
	Label   string            `yaml:"label"`
	Type    TypeItem          `yaml:"type"`
	Value   string            `yaml:"value"`
	Options []string          `yaml:"options,flow"`
	Opts    map[string]string `yaml:"optobj,flow"`
	cfg     *Cfg
}

func (i *Item) defaults(c *Cfg) {
	i.cfg = c
	if i.Label == "" {
		i.Label = i.Name
	}
}

func (i *Item) getfs(we *widget.Entry) {
	switch i.Type {
	case File:
		lfo := func(uri fyne.URIReadCloser, err error) {
			if err == nil {
				val := uri.URI().Path()
				we.SetText(val)
			}
		}
		dlg := dialog.NewFileOpen(lfo, i.cfg.win)
		du, e := storage.ListerForURI(storage.NewFileURI(i.Value))
		if e == nil {
			dlg.SetLocation(du)
		}
		dlg.Show()

	// func NewFileSave
	// func NewFileSave(callback func(fyne.URIWriteCloser, error), parent fyne.Window) *FileDialog
	// NewFileSave creates a file dialog allowing the user to choose a file to save to (new or overwrite). If the user chooses an existing file they will be asked if they are sure. The dialog will appear over the window specified when Show() is called.

	case Dir:
		lfo := func(uri fyne.ListableURI, err error) {
			if err == nil {
				val := uri.Path()
				we.SetText(val)
			}
		}
		dlg := dialog.NewFolderOpen(lfo, i.cfg.win)
		du, e := storage.ListerForURI(storage.NewFileURI(i.Value))
		if e == nil {
			dlg.SetLocation(du)
		}
		dlg.Show()
	}
}

func (i *Item) changed(val string) {
	i.Value = val
}

func (i *Item) widgets() (fyne.CanvasObject, fyne.CanvasObject) {
	label := widget.NewLabel(i.Label)
	switch i.Type {
	case Text:
		w := widget.NewEntry()
		w.OnChanged = func(val string) { i.changed(val) }
		w.SetText(i.Value)
		return label, w
	case Password:
		w := widget.NewPasswordEntry()
		w.OnChanged = func(val string) { i.changed(val) }
		w.SetText(i.Value)
		return label, w
	case Select:
		s := widget.NewSelect(i.Options, func(val string) { i.changed(val) })
		s.SetSelected(i.Value)
		return label, s
	case File, Dir:
		w := widget.NewEntry()
		w.SetText(i.Value)
		w.OnChanged = func(val string) { i.changed(val) }
		b := widget.NewButton("...", func() { i.getfs(w) })
		l := container.NewHBox(b, w)
		// l.Resize()
		return label, l
		// Execute TypeItem = "execute"
		// Cancel  TypeItem = "cancel"
		// Next    TypeItem = "next"
		// Back    TypeItem = "back"
	default:
		return label, label
	}
}
