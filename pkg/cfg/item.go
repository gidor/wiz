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
	"bytes"
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	extdlg "github.com/gidor/wiz/pkg/dlg"
	"github.com/gidor/wiz/pkg/runner"
)

var (
	allvalues map[string]string
	run       *runner.Runner
)

func init() {
	allvalues = make(map[string]string)

}

// A Form item
type Item struct {

	// name: text
	Name    string   `yaml:"name"`
	Label   string   `yaml:"label"`
	Type    TypeItem `yaml:"type"`
	value   string   `yaml:"value"`
	Options []string `yaml:"options,flow"`
	Todo    Action   `yaml:"action"`
	cfg     *Cfg
}

func (i *Item) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var dy struct {
		// name: text
		Name    string   `yaml:"name"`
		Label   string   `yaml:"label"`
		Type    string   `yaml:"type"`
		Value   string   `yaml:"value"`
		Todo    Action   `yaml:"action"`
		Options []string `yaml:"options,flow"`
	}
	if err := unmarshal(&dy); err != nil {
		fmt.Println(err.Error())

		return err
	}
	i.Todo = dy.Todo
	i.value = dy.Value
	i.Name = dy.Name
	i.Label = dy.Label
	// i.Dynamic = dy.Dynamic
	i.Options = dy.Options
	switch dy.Type {
	case string(Text):
		i.Type = Text
	case string(Password):
		i.Type = Password
	case string(File):
		i.Type = File
	case string(FileOpen):
		i.Type = FileOpen
	case string(FileSave):
		i.Type = FileSave
	case string(Dir):
		i.Type = Dir
	case string(Select):
		i.Type = Select
	case string(Execute):
		i.Type = Execute
		if i.Todo.Title == "" {
			i.Todo.Title = strings.Title(dy.Type)
		}
	case string(Cancel):
		i.Type = Cancel
		if i.Todo.Title == "" {
			i.Todo.Title = strings.Title(dy.Type)
		}
	case string(Next):
		i.Type = Next
		if i.Todo.Title == "" {
			i.Todo.Title = strings.Title(dy.Type)
		}
	case string(Back):
		i.Type = Back
		if i.Todo.Title == "" {
			i.Todo.Title = strings.Title(dy.Type)
		}
	default:
		i.Type = Text
	}
	if i.Label == "" {
		i.Label = i.Name
	}
	return nil
}

func (i *Item) defaults(c *Cfg) {
	i.cfg = c
}

// run task or navigate beetween forms
func (i *Item) execute() {
	if i.Todo.Execute != "" {
		if run == nil {
			r, err := runner.NewRunner(i.cfg.Taskfile())
			if err != nil {
				fmt.Println(err.Error())
				return
			} else {
				run = r
			}
		}
		run.Run(i.Todo.Execute, allvalues)
		return
	}
	if i.Todo.Goto != "" {
		for _, form := range i.cfg.Panels {
			if form.Title == i.Todo.Goto {
				form.render()
				break
			}
		}
		return
	}
}

// run task an get output as sing slice
func (i *Item) executeTo() []string {
	if i.Todo.Execute != "" {
		if run == nil {
			r, err := runner.NewRunner(i.cfg.Taskfile())
			if err != nil {
				fmt.Println(err.Error())
				return nil
			} else {
				run = r
			}
		}

		var buffer bytes.Buffer
		run.RunOn(i.Todo.Execute, &buffer, allvalues)
		return strings.Split(buffer.String(), "\n")
	}
	return nil
}

// open and manage dialog in order to navigate the file system
func (i *Item) getfs(we *widget.Entry) {
	switch i.Type {
	case File, FileOpen:
		lfo := func(uri fyne.URIReadCloser, err error) {
			if err == nil {
				if uri != nil {
					val := uri.URI().Path()
					// val := uri.URI().String()
					we.SetText(val)
				}
			}
		}
		dlg := dialog.NewFileOpen(lfo, i.cfg.win)
		du, e := storage.ListerForURI(storage.NewFileURI(i.Val()))
		if e == nil {
			dlg.SetLocation(du)
		}
		dlg.Show()
	case FileSave:
		lfo := func(uri fyne.URIWriteCloser, err error) {
			if err == nil {
				if uri != nil {
					val := uri.URI().Path()
					// val := uri.URI().String()
					we.SetText(val)
				}
			}
		}
		dlg := dialog.NewFileSave(lfo, i.cfg.win)
		du, e := storage.ListerForURI(storage.NewFileURI(i.Val()))
		if e == nil {
			dlg.SetLocation(du)
		}
		dlg.Show()
	case Dir:
		lfo := func(uri fyne.ListableURI, err error) {
			if err == nil {
				if uri != nil {
					val := uri.Path()
					// val := uri.String()
					we.SetText(val)
				}
			}
		}
		dlg := dialog.NewFolderOpen(lfo, i.cfg.win)
		du, e := storage.ListerForURI(storage.NewFileURI(i.Val()))
		if e == nil {
			dlg.SetLocation(du)
		}
		dlg.Show()
	}
}

// open and manage dialog in order to navigate the file system
func (i *Item) getfs_ext(we *widget.Entry) {
	switch i.Type {
	case File, FileOpen:
		file, err := extdlg.DlgFile(false, i.Label, i.Val(), i.Options)
		if err == nil {
			we.SetText(file)
		}
	case FileSave:
		file, err := extdlg.DlgFile(true, i.Label, i.Val(), i.Options)
		if err == nil {
			we.SetText(file)
		}
	case Dir:
		file, err := extdlg.DlgDir(i.Label, i.Val())
		if err == nil {
			we.SetText(file)
		}
	}
}

// get the item value from default or from the runtime cache
func (i *Item) Val() string {
	val, ok := allvalues[i.Name]
	if ok {
		return val
	} else {
		allvalues[i.Name] = i.value
		return i.value
	}
}

// set new val for the from item in the runtime cache
func (i *Item) changed(val string) {
	allvalues[i.Name] = val
}

//get the widgets for the items to be rendered in a form layout row
func (i *Item) widgets() (fyne.CanvasObject, fyne.CanvasObject) {
	label := widget.NewLabel(i.Label)
	switch i.Type {
	case Text:
		w := widget.NewEntry()
		w.OnChanged = func(val string) { i.changed(val) }
		w.SetText(i.Val())
		return label, w
	case Password:
		w := widget.NewPasswordEntry()
		w.OnChanged = func(val string) { i.changed(val) }
		w.SetText(i.Val())
		return label, w
	case Select:
		if i.Options == nil || len(i.Options) == 0 {
			i.Options = i.executeTo()
			if i.Options == nil {
				i.Options = make([]string, 0)
			}
		}
		s := widget.NewSelect(i.Options, func(val string) { i.changed(val) })
		s.SetSelected(i.Val())
		return label, s
	case File, Dir, FileOpen, FileSave:
		w := widget.NewEntry()
		w.SetText(i.Val())
		w.OnChanged = func(val string) { i.changed(val) }
		b := widget.NewButton("...", func() { i.getfs_ext(w) })
		// l := container.NewHBox(b, w)
		l := container.New(layout.NewFormLayout(), b, w)
		// l.Resize()
		return label, l
	case Execute, Cancel, Next, Back:
		w := widget.NewButton(i.Todo.Title, i.execute)
		return label, w
		// Cancel  TypeItem = "cancel"
		// Next    TypeItem = "next"
		// Back    TypeItem = "back"
	default:
		return label, label
	}
}
