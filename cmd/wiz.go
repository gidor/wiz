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

package cmd

import (
	"fmt"
	// "log"
	"os"
	"path/filepath"

	"github.com/gidor/wiz/pkg/cfg"
	log "github.com/gidor/wiz/pkg/logwrapper"

	"github.com/spf13/pflag"
	// "github.com/go-task/task/v3"
)

const (
	synopsis = `Usage: wiz [-v] [-f config]
flag 

default  configuration wiz.yaml
Options:
`
	version = "0.1.beta"
)

func setuplog(dir string) {

	logp := filepath.Join(dir, "wiz.log")
	logfile, err := os.OpenFile(logp, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(logfile)

}

func Start() {

	var (
		showversion bool
		showhelp    bool
		dry         bool
		dir         string
		configpath  string
	)

	pflag.Usage = func() {
		fmt.Print(synopsis)
		pflag.PrintDefaults()
	}

	pflag.BoolVarP(&showversion, "version", "v", false, "show version")
	pflag.BoolVar(&dry, "dry", false, "dso notexecute any task")
	pflag.StringVarP(&dir, "dir", "d", "", "sets directory of execution")
	pflag.StringVarP(&configpath, "configfile", "f", "", `choose config file. Defaults to "wiz.yaml"`)
	pflag.BoolVarP(&showhelp, "help", "h", false, "show help")

	pflag.Parse()

	if showhelp {
		pflag.Usage()
	}
	if dir != "" && configpath != "" {
		log.Fatal("got dir and confifgile")
	} else if configpath != "" {
		dir = filepath.Dir(configpath)
		configpath = filepath.Base(configpath)
	} else {
		configpath = "wiz.yaml"
	}
	if dir == "" {
		// d, err := os.Executable()
		d, err := os.Getwd()
		if err != nil {
			log.Println(err)
		}
		// dir = filepath.Dir(d)
		dir = d
	}
	setuplog(dir)

	cfgp := filepath.Join(dir, configpath)

	if _, err := os.Stat(cfgp); err != nil {
		log.Fatal(err)
	} else {
		cfg := cfg.GetCfg(cfgp)
		cfg.Start()
	}
}
