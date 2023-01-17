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
	dir          string
	entrypoint   string
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
	output := taskfile.Output{}

	//interleaved|group|prefixed]")
	output.Name = "interleaved"
	//, "output-group-begin", "", "message template to print before a task's grouped output")
	output.Group.Begin = ""
	//, "output-group-end", "", "message template to print after a task's grouped output")
	output.Group.End = ""

	r.executor = task.Executor{
		Force:       false,
		Watch:       false,
		Verbose:     false,
		Silent:      true,
		Dir:         r.dir,
		Dry:         r.dryrun,
		Entrypoint:  r.entrypoint,
		Summary:     false,
		Parallel:    false,
		Color:       false,
		Concurrency: 0,

		Stdin:       os.Stdin,
		Stdout:      os.Stdout,
		Stderr:      os.Stderr,
		OutputStyle: output,
	}

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

func (r *Runner) Run(task string, params map[string]string) error {
	if r.reloadPolicy == Always {
		if !r.reload() {
			return nil
		}
	}
	r.show(task, params)
	if r.dryrun {
		return nil
	}
	exe := r.executor
	call := taskfile.Call{Task: task, Vars: &taskfile.Vars{}}
	for k, v := range params {
		call.Vars.Set(k, taskfile.Var{Static: v})
	}
	ctx := context.Background()
	if err := exe.Run(ctx, call); err != nil {
		log.Print(err)
		return err
	}
	return nil

}

func (r *Runner) RunOn(task string, output io.Writer, params map[string]string) error {
	e := r.executor
	if r.reloadPolicy == Always {
		log.Print("Reloading")
		if !r.reload() {
			return nil
		}
	}
	r.show(task, params)
	if r.dryrun {
		return nil
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
		return err
	}
	e.Stdout = legacyout
	return nil

}
