package appname

import (
	"os"
	"path/filepath"
)

func GetAPPName() (app string) {
	app = "APP_Name"
	exe, err := os.Executable()
	if err == nil {
		app = filepath.Base(exe)
	}

	return
}
