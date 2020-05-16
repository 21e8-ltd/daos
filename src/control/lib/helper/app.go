//
// (C) Copyright 2020 Intel Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// GOVERNMENT LICENSE RIGHTS-OPEN SOURCE SOFTWARE
// The Government's rights to use, modify, reproduce, release, perform, display,
// or disclose this software are subject to the terms of the Apache License as
// provided in Contract No. 8F-30005.
// Any reproduction of computer software, computer software documentation, or
// portions thereof marked with this legend must also reproduce the markings.
//

package helper

import (
	"fmt"
	// "io"
	// "io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/daos-stack/daos/src/control/build"
	// "github.com/daos-stack/daos/src/control/common"
	"github.com/daos-stack/daos/src/control/logging"
	"github.com/daos-stack/daos/src/control/pbin"
)

// App is a framework for an external helper application to be invoked by one or
// more DAOS processes.
type App struct {
	log            logging.Logger
	allowedCallers []string
}

// Name returns the name of the application binary.
func (a *App) Name() string {
	return filepath.Base(os.Args[0])
}

// ParentProcessName returns the name of the binary that invoked this application, or an error.
func (a *App) ParentProcessName() (string, error) {
	pPath, err := os.Readlink(fmt.Sprintf("/proc/%d/exe", os.Getppid()))
	if err != nil {
		return "", a.logError(errors.Wrap(err, "failed to identify parent process binary"))
	}

	return filepath.Base(pPath), nil
}

// logError is a convenience method that logs an error and returns the same error
func (a *App) logError(err error) error {
	a.log.Error(err.Error())
	return err
}

// Run executes the helper application process.
func (a *App) Run() error {
	return nil
}

func (a *App) checkParentName() error {
	parentName, err := a.ParentProcessName()
	if err != nil {
		return err
	}
	if !isCallerPermitted(parentName) {
		return a.logError(errors.Errorf("%s (version %s) may only be invoked by: %v",
			a.Name(), build.DaosVersion, a.allowedCallers))
	}

	return nil
}

func (a *App) isCallerPermitted(callerName string) bool {
	for _, name := range a.allowedCallers {
		if callerName == name {
			return true
		}
	}
	return false
}
