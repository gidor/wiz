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

package runner

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/go-task/task/v3"
	"github.com/go-task/task/v3/taskfile"
)

type ReloadPolicy int
type Verbosity int

const (
	Never   ReloadPolicy = 0
	Always  ReloadPolicy = 1
	Silent  Verbosity    = 0
	Verbose Verbosity    = 1
)

type Runner struct {
	reloadPolicy ReloadPolicy
	verbose      Verbosity
	dryrun       bool
	taskpath     string
	executor     task.Executor
	// taskfile taskfile.Taskfile
}

func NewRunner(taskpath string, reload ReloadPolicy, verbose Verbosity, dry bool) (*Runner, error) {
	dir := filepath.Dir(taskpath)
	entrypoint := filepath.Base(taskpath)

	executor := task.Executor{
		Force:       false,
		Watch:       false,
		Verbose:     false,
		Silent:      true,
		Dir:         dir,
		Dry:         false,
		Entrypoint:  entrypoint,
		Summary:     false,
		Parallel:    false,
		Color:       false,
		Concurrency: 0,

		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,

		OutputStyle: "", // "interleaved",
	}

	if err := executor.Setup(); err != nil {
		log.Print("error in setup ", taskpath, err)
		//executor.Logger.Fatal(err)
		return nil, err
	}

	re := Runner{
		dryrun:       dry,
		reloadPolicy: reload,
		verbose:      verbose,
		taskpath:     taskpath,
		executor:     executor,
		// taskfile: nil,
	}

	return &re, nil
}

func (r *Runner) reload() bool {
	// dir := filepath.Dir(r.taskpath)
	// entrypoint := filepath.Base(r.taskpath)

	// r.executor = task.Executor{
	// 	Force:       false,
	// 	Watch:       false,
	// 	Verbose:     false,
	// 	Silent:      true,
	// 	Dir:         dir,
	// 	Dry:         false,
	// 	Entrypoint:  entrypoint,
	// 	Summary:     false,
	// 	Parallel:    false,
	// 	Color:       false,
	// 	Concurrency: 0,

	// 	Stdin:  os.Stdin,
	// 	Stdout: os.Stdout,
	// 	Stderr: os.Stderr,

	// 	OutputStyle: "", // "interleaved",
	// }

	if err := r.executor.Setup(); err != nil {
		log.Print("error in setup ", r.taskpath, err)
		return false
	}
	return true
}

func (r *Runner) show(task string, params map[string]string) {
	if r.verbose == Verbose {
		fmt.Println("taskfile: ", r.taskpath)
		fmt.Println("task:", task)
		fmt.Println("parmas:")
		for key, value := range params {
			fmt.Println("key:", key, "Value:", value)
		}
	}

}

func (r *Runner) Run(task string, params map[string]string) {
	if r.reloadPolicy == Always {
		if !r.reload() {
			return
		}
	}
	r.show(task, params)
	if r.dryrun {
		return
	}
	exe := r.executor
	call := taskfile.Call{Task: task, Vars: &taskfile.Vars{}}
	for k, v := range params {
		call.Vars.Set(k, taskfile.Var{Static: v})
	}
	ctx := context.Background()
	if err := exe.Run(ctx, call); err != nil {
		log.Print("Error runing ", task, err)
	}

}

func (r *Runner) RunOn(task string, output io.Writer, params map[string]string) {
	e := r.executor
	if r.reloadPolicy == Always {
		if !r.reload() {
			return
		}
	}
	r.show(task, params)
	if r.dryrun {
		return
	}
	legacyout := e.Stdout
	e.Stdout = output
	call := taskfile.Call{Task: task, Vars: &taskfile.Vars{}}
	for k, v := range params {
		call.Vars.Set(k, taskfile.Var{Static: v})
	}
	ctx := context.Background()
	if err := e.Run(ctx, call); err != nil {
		log.Print(err)
	}
	e.Stdout = legacyout

}
