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

package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"github.com/gidor/wiz/pkg/cfg"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

var (
	wizApp  fyne.App
	wizWin  fyne.Window
	wizTabs container.AppTabs
)

func Setup(c *cfg.Cfg) {
	wizApp = app.New()
	wizWin = wizApp.NewWindow("Hello")
	//.SetContent(widget.NewLabel("Hello"))

}

func Start() {
	wizWin.Show()
	wizApp.Run()
	tidyUp()
}
func tidyUp() {
	fmt.Println("Exited")
}
