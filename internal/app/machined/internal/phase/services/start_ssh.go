// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package services

import (
	"github.com/talos-systems/talos/internal/app/machined/internal/phase"
	"github.com/talos-systems/talos/internal/app/machined/pkg/system"
	"github.com/talos-systems/talos/internal/app/machined/pkg/system/services"
	"github.com/talos-systems/talos/internal/pkg/runtime"
)

// StartSSH represents the task to start system containerd.
type StartSSH struct{}

// NewStartSSHTask initializes and returns an Services task.
func NewStartSSHTask() phase.Task {
	return &StartSSH{}
}

// TaskFunc returns the runtime function.
func (task *StartSSH) TaskFunc(mode runtime.Mode) phase.TaskFunc {
	return task.standard
}

func (task *StartSSH) standard(r runtime.Runtime) (err error) {
	system.Services(r.Config()).LoadAndStart(&services.SSH{})

	return nil
}
