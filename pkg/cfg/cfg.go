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
	"io/ioutil"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	"fyne.io/fyne/v2/widget"
	goyaml "gopkg.in/yaml.v3"
)

type Cfg struct {
	Height  float32
	Width   float32
	Title   string  `yaml:"title"`
	Panels  []*Form `yaml:",flow"`
	lasterr error
	app     fyne.App
	win     fyne.Window
}

func (c *Cfg) Start() {
	c.Size()
	c.win.Show()
	c.app.Run()
	tidyUp()
}

func max(a float32, b float32) float32 {
	if a > b {
		return a
	}
	return b
}
func (c *Cfg) Size() {
	size := c.win.Content().MinSize()
	c.win.CenterOnScreen()
	c.win.Resize(fyne.NewSize(max(c.Width, size.Width), max(c.Height, size.Height)))
}

func tidyUp() {
	fmt.Println("Exited")
}

func (c *Cfg) defaults() {
	if c.Height == 0 {
		c.Height = 300
	}
	if c.Width == 0 {
		c.Width = 400
	}
	c.app = app.New()
	c.win = c.app.NewWindow(c.Title)
	items := make([]*fyne.MenuItem, 0, len(c.Panels)+2)
	for _, form := range c.Panels {
		form.defaults(c)
		item, err := form.MenuItem()
		if err == nil {
			items = append(items, item)
		}
	}

	//.SetContent(widget.NewLabel("Hello"))
	items = append(items, fyne.NewMenuItemSeparator())
	quit := fyne.NewMenuItem("Exit", func() { c.app.Quit() })
	quit.IsQuit = true
	items = append(items, quit)
	todomenu := fyne.NewMenu("Todo", items...)
	c.win.SetMainMenu(fyne.NewMainMenu(todomenu))
	c.win.SetContent(widget.NewLabel("TODO"))

}

func GetCfg(path string) Cfg {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	cfg := Cfg{}

	err = goyaml.Unmarshal(b, &cfg)
	if err != nil {
		cfg.lasterr = err
		panic(err)
	}
	cfg.defaults()
	return cfg
}
