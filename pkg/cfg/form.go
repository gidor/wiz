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
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	// "fyne.io/fyne/v2/widget"
)

type Action struct {
	Execute string `yaml:"execute"`
}

// action:
//   execute: task {{.dir}} {{.select}} {{.text}} {{.dir}} {{.file}}
// form:

type Form struct {
	Title   string  `yaml:"title"`
	Visible string  `yaml:"public"`
	Todo    Action  `yaml:"action"`
	Items   []*Item `yaml:"form,flow"`
	public  bool
	cfg     *Cfg
}

// get the menuitem that will diplay the panell
func (f *Form) MenuItem() (*fyne.MenuItem, error) {
	if f.HasMenu() {
		mi := fyne.NewMenuItem(f.Title, f.render)
		return mi, nil
	} else {
		return nil, fmt.Errorf("forms %s is private", f.Title)
	}
}

// rendere the form
func (f *Form) render() {
	objs := make([]fyne.CanvasObject, 0, (len(f.Items)+1)*2)
	// collect form widget
	for _, item := range f.Items {
		label, content := item.widgets()
		objs = append(objs, label, content)
	}
	form := container.New(layout.NewFormLayout(), objs...)
	f.cfg.win.SetContent(form)

}

// execute the form
func (f *Form) execute() {

}

func (f *Form) HasMenu() bool {
	return f.public
}

func (f *Form) defaults(c *Cfg) {
	f.cfg = c
	switch strings.ToLower(f.Visible) {
	case "yes":
		f.public = true
	case "no":
		f.public = false
	default:
		f.public = true
	}
	for _, item := range f.Items {
		item.defaults(c)

	}
}
