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

package dlg

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/sqweek/dialog"
)

// dlg_ask for file using dialog

func DlgFile(save bool, title string, path string, exts []string) (string, error) {
	dlg := dialog.File()
	dlg.Title(title)

	for _, ext := range exts {
		dlg.Filter("Files "+ext, strings.TrimLeft(ext, "."))
	}

	dlg.Filter("All Files", "*")

	if path == "" || path == "." {
		dlg.StartDir, _ = os.Getwd()
	} else {
		dlg.StartDir, _ = filepath.Abs(filepath.Dir(path))
	}
	if save {
		return dlg.Save()
		// return file, err
	} else {
		return dlg.Load()
		// file, err := dlg.Load()
		// return file, err
	}
}

func DlgDir(title string, path string) (string, error) {
	if path == "" {
		path = "."
	}
	dlg := dialog.Directory()
	dlg.StartDir = path
	dlg.Title(title)
	file, err := dlg.Browse()
	return file, err
}
