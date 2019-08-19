/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package network

import (
	"log"

	"github.com/talos-systems/talos/internal/app/machined/internal/phase"
	"github.com/talos-systems/talos/internal/app/machined/internal/phase/rootfs/etc"
	"github.com/talos-systems/talos/internal/app/machined/internal/platform"
	"github.com/talos-systems/talos/internal/app/machined/internal/runtime"
	"github.com/talos-systems/talos/internal/app/networkd/pkg/networkd"
	"github.com/talos-systems/talos/internal/app/networkd/pkg/nic"
	"github.com/talos-systems/talos/internal/pkg/kernel"
	"github.com/talos-systems/talos/pkg/constants"
	"github.com/talos-systems/talos/pkg/userdata"
)

// UserDefinedNetwork represents the UserDefinedNetwork task.
type UserDefinedNetwork struct{}

// NewUserDefinedNetworkTask initializes and returns an UserDefinedNetwork task.
func NewUserDefinedNetworkTask() phase.Task {
	return &UserDefinedNetwork{}
}

// RuntimeFunc returns the runtime function.
func (task *UserDefinedNetwork) RuntimeFunc(mode runtime.Mode) phase.RuntimeFunc {
	switch mode {
	case runtime.Standard:
		return task.runtime
	default:
		return nil
	}
}

// nolint: gocyclo
func (task *UserDefinedNetwork) runtime(platform platform.Platform, data *userdata.UserData) (err error) {
	nwd, err := networkd.New()
	if err != nil {
		return err
	}

	// Convert links to nic
	log.Println("Discovering local network interfaces")
	netconf, err := nwd.Discover()
	if err != nil {
		return err
	}

	// Configure specified interface
	netIfaces := make([]*nic.NetworkInterface, 0, len(netconf))
	var iface *nic.NetworkInterface
	for link, opts := range netconf {
		log.Printf("Creating interface %s", link.Name)
		iface, err = nic.Create(opts...)
		if err != nil {
			return err
		}

		netIfaces = append(netIfaces, iface)
	}

	// kick off the addressing mechanism
	// Add any necessary routes
	log.Println("Configuring interface addressing")
	if err = nwd.Configure(netIfaces...); err != nil {
		return err
	}

	// Create /etc/hosts.
	// Priority is:
	// 1. Userdata
	// 2. Kernel Arg
	// 3. DHCP response
	// 4. failsafe - talos-<ip addr>
	var hostname string
	kernelHostname := kernel.ProcCmdline().Get(constants.KernelParamHostname).First()
	switch {
	case data.Networking != nil && data.Networking.OS != nil && data.Networking.OS.Hostname != "":
		hostname = data.Networking.OS.Hostname
	case kernelHostname != nil:
		hostname = *kernelHostname
	case nwd.Hostname() != "":
		hostname = nwd.Hostname(netIfaces...)
	default:
		// will default in Hosts()
	}

	return etc.Hosts(hostname)
}
