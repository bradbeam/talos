/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package userdata

import (
	"github.com/talos-systems/talos/internal/app/machined/internal/phase"
	"github.com/talos-systems/talos/internal/app/machined/internal/platform"
	"github.com/talos-systems/talos/internal/app/machined/internal/runtime"
	"github.com/talos-systems/talos/pkg/userdata"

	kubeadmapi "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm"
)

// UserData represents the UserData task.
type UserData struct{}

// NewUserDataTask initializes and returns an UserData task.
func NewUserDataTask() phase.Task {
	return &UserData{}
}

// RuntimeFunc returns the runtime function.
func (task *UserData) RuntimeFunc(mode runtime.Mode) phase.RuntimeFunc {
	switch mode {
	case runtime.Standard:
		return task.standard
	case runtime.Container:
		return task.container
	default:
		return nil
	}
}

func (task *UserData) standard(platform platform.Platform, data *userdata.UserData) (err error) {
	var d *userdata.UserData
	d, err = platform.UserData()
	if err != nil {
		return err
	}
	*data = *d

	return nil
}

func (task *UserData) container(platform platform.Platform, data *userdata.UserData) (err error) {
	var d *userdata.UserData
	d, err = platform.UserData()
	if err != nil {
		return err
	}
	*data = *d

	data.Services.Kubeadm.IgnorePreflightErrors = []string{"FileContent--proc-sys-net-bridge-bridge-nf-call-iptables", "Swap", "SystemVerification"}
	initConfiguration, ok := data.Services.Kubeadm.Configuration.(*kubeadmapi.InitConfiguration)
	if ok {
		initConfiguration.ClusterConfiguration.ComponentConfigs.Kubelet.FailSwapOn = false
		// See https://github.com/kubernetes/kubernetes/issues/58610#issuecomment-359552443
		maxPerCore := int32(0)
		initConfiguration.ClusterConfiguration.ComponentConfigs.KubeProxy.Conntrack.MaxPerCore = &maxPerCore
	}

	return nil
}
