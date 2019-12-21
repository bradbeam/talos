/* This Source Code Form is subject to the terms of the Mozilla Public
* License, v. 2.0. If a copy of the MPL was not distributed with this
* file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package services

import (
	"context"
	"fmt"

	containerdapi "github.com/containerd/containerd"
	"github.com/containerd/containerd/oci"
	specs "github.com/opencontainers/runtime-spec/specs-go"

	"github.com/talos-systems/talos/internal/app/machined/pkg/system/runner"
	"github.com/talos-systems/talos/internal/app/machined/pkg/system/runner/containerd"
	"github.com/talos-systems/talos/internal/app/machined/pkg/system/runner/restart"
	"github.com/talos-systems/talos/internal/pkg/conditions"
	"github.com/talos-systems/talos/internal/pkg/runtime"
	"github.com/talos-systems/talos/pkg/constants"
)

// SSH implements the Service interface. It serves as the concrete type with
// the required methods.
type SSH struct{}

// ID implements the Service interface.
func (k *SSH) ID(config runtime.Configurator) string {
	return "ssh"
}

// PreFunc implements the Service interface.
func (k *SSH) PreFunc(ctx context.Context, config runtime.Configurator) error {
	importer := containerd.NewImporter(constants.SystemContainerdNamespace, containerd.WithContainerdAddress(constants.SystemContainerdAddress))

	return importer.Import(&containerd.ImportRequest{
		Path: "/usr/images/ssh.tar",
		Options: []containerdapi.ImportOpt{
			containerdapi.WithIndexName("talos/ssh"),
		},
	})
}

// PostFunc implements the Service interface.
func (k *SSH) PostFunc(config runtime.Configurator) (err error) {
	return nil
}

// Condition implements the Service interface.
func (k *SSH) Condition(config runtime.Configurator) conditions.Condition {
	return conditions.None()
}

// DependsOn implements the Service interface.
func (k *SSH) DependsOn(config runtime.Configurator) []string {
	return []string{"system-containerd"}
}

// Runner implements the Service interface.
func (k *SSH) Runner(config runtime.Configurator) (runner.Runner, error) {
	image := "talos/ssh"
	// Set the process arguments.
	args := runner.Args{
		ID: k.ID(config),
		ProcessArgs: []string{
			"/usr/sbin/sshd",
			"-D",
		},
	}
	// Set the required SSH mounts.
	mounts := []specs.Mount{
		{Type: "bind", Destination: "/dev", Source: "/dev", Options: []string{"rbind", "rshared", "rw"}},
		{Type: "sysfs", Destination: "/sys", Source: "/sys", Options: []string{"bind", "ro"}},
		{Type: "cgroup", Destination: "/sys/fs/cgroup", Options: []string{"rbind", "rshared", "rw"}},
		{Type: "bind", Destination: "/rootfs", Source: "/", Options: []string{"rbind", "rshared", "rw"}},
		{Type: "bind", Destination: "/lib/modules", Source: "/lib/modules", Options: []string{"bind", "ro"}},
	}
	env := []string{}
	for key, val := range config.Machine().Env() {
		env = append(env, fmt.Sprintf("%s=%s", key, val))
	}
	return restart.New(containerd.NewRunner(
		config.Debug(),
		&args,
		runner.WithContainerdAddress(constants.SystemContainerdAddress),
		runner.WithContainerImage(image),
		runner.WithEnv(env),
		runner.WithOCISpecOpts(
			containerd.WithRootfsPropagation("shared"),
			oci.WithMounts(mounts),
			oci.WithHostNamespace(specs.PIDNamespace),
			oci.WithHostNamespace(specs.NetworkNamespace),
			oci.WithParentCgroupDevices,
			oci.WithPrivileged,
		),
	),
		restart.WithType(restart.Forever),
	), nil
}

/*
func (k *SSH) HealthFunc(runtime.Configurator) health.Check {
	return nil
}

func (k *SSH) HealthSettings(runtime.Configurator) *health.Settings {
	return &health.DefaultSettings
}
*/
