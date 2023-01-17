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
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/gidor/wiz/pkg/runner"
	goyaml "gopkg.in/yaml.v3"
)

type Dsize struct {
	W float32
	H float32
}

type Cfg struct {
	Menu           string  `yaml:"menu"`
	Msg            string  `yaml:"msg"`
	Title          string  `yaml:"title"`
	Panels         []*Form `yaml:",flow"`
	taskentrypoint string
	entrypoint     string
	lasterr        error
	height         float32
	width          float32
	app            fyne.App
	win            fyne.Window
	dry            bool
	verbose        runner.Verbosity
	reload         runner.ReloadPolicy
	popUp          *widget.PopUp
}

func (c *Cfg) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var dy struct {
		Taskentrypoint string  `yaml:"taskfile"`
		Menu           string  `yaml:"menu"`
		Minsize        Dsize   `yaml:"minsize"`
		Msg            string  `yaml:"msg"`
		Title          string  `yaml:"title"`
		Panels         []*Form `yaml:",flow"`
	}

	if err := unmarshal(&dy); err != nil {
		c.lasterr = err
		fmt.Println(err.Error())
		return err
	}
	if dy.Minsize.H > 0 {
		c.height = dy.Minsize.H
	} else {
		c.height = 300
	}
	if dy.Minsize.W > 0 {
		c.width = dy.Minsize.W
	} else {
		c.width = 400
	}

	c.Msg = dy.Msg
	c.Menu = dy.Menu
	if c.Menu == "" {
		c.Menu = "Todo"
	}
	c.taskentrypoint = dy.Taskentrypoint
	c.Title = dy.Title
	c.Panels = dy.Panels
	return nil
}

func (c *Cfg) Taskfile() string {
	if filepath.IsAbs(c.taskentrypoint) {
		return c.taskentrypoint
	} else {
		return filepath.Join(filepath.Dir(c.entrypoint), c.taskentrypoint)
	}

}

// start showing the Configured wizard
func (c *Cfg) Start(verbose bool, dry bool, reload bool) {
	if verbose {
		c.verbose = runner.Verbose
	} else {
		c.verbose = runner.Silent
	}
	c.dry = dry
	if reload {
		c.reload = runner.Always
	} else {
		c.reload = runner.Never
	}
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

// set window size
func (c *Cfg) Size() {
	size := c.win.Content().MinSize()
	c.win.CenterOnScreen()
	c.win.Resize(fyne.NewSize(max(c.width, size.Width), max(c.height, size.Height)))
}

func tidyUp() {
	fmt.Println("Exited")
}

func (c *Cfg) ShowRunning() {
	if c.win != nil {
		c.popUp = widget.NewModalPopUp(
			container.NewVBox(
				widget.NewLabel("Running"),
			),
			c.win.Canvas(),
		)
		c.popUp.Show()
	}
}

func (c *Cfg) ShowFinished(e error) {
	if c.win != nil {
		if c.popUp != nil {
			if !c.popUp.Hidden {
				c.popUp.Hide()
			}
			c.popUp = nil
		}
		var msg string
		var modal *widget.PopUp
		if e == nil {
			msg = "Done"
		} else {
			msg = e.Error()
		}

		modal = widget.NewModalPopUp(
			container.NewVBox(
				widget.NewLabel(msg),
				widget.NewButton("OK", func() { modal.Hide() }),
			),
			c.win.Canvas(),
		)
		modal.Show()
	}
}

// render the form
func (c *Cfg) render() {
	c.win.SetContent(widget.NewLabel(c.Msg))
}

func (c *Cfg) defaults() {
	c.app = app.New()
	c.win = c.app.NewWindow(c.Title)
	items := make([]*fyne.MenuItem, 0, len(c.Panels)+2)
	items = append(items, fyne.NewMenuItem("Home", c.render))
	items = append(items, fyne.NewMenuItemSeparator())
	for _, form := range c.Panels {
		form.defaults(c)
		item, err := form.MenuItem()
		if err == nil {
			items = append(items, item)
		} else {
			c.lasterr = err
		}
	}

	items = append(items, fyne.NewMenuItemSeparator())
	quit := fyne.NewMenuItem("Exit", func() { c.app.Quit() })
	quit.IsQuit = true
	items = append(items, quit)
	todomenu := fyne.NewMenu(c.Menu, items...)
	c.win.SetMainMenu(fyne.NewMainMenu(todomenu))
	c.render()

}

//  Configuration  factory
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
	cfg.entrypoint = path
	cfg.defaults()
	return cfg
}
