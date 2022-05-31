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
	dlg.Filter("All Files", "*")

	for _, ext := range exts {
		dlg.Filter("FILES "+ext, strings.TrimLeft(ext, "."))
	}

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
