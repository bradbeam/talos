/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package upgrade

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/talos-systems/talos/internal/pkg/blockdevice/bootloader/syslinux"
	"github.com/talos-systems/talos/internal/pkg/blockdevice/probe"
	"github.com/talos-systems/talos/internal/pkg/constants"
	"github.com/talos-systems/talos/internal/pkg/install"
	"github.com/talos-systems/talos/internal/pkg/kernel"
	"github.com/talos-systems/talos/internal/pkg/kubernetes"
	"github.com/talos-systems/talos/internal/pkg/mount"
	"golang.org/x/sys/unix"
)

// NewUpgrade initiates a Talos upgrade
func NewUpgrade(url string) (err error) {
	if err = upgradeRoot(url); err != nil {
		return
	}
	if err = upgradeBoot(url); err != nil {
		return
	}

	// cordon/drain
	var kubeHelper *kubernetes.Helper
	if kubeHelper, err = kubernetes.NewHelper(); err != nil {
		return
	}

	var hostname string
	if hostname, err = os.Hostname(); err != nil {
		return
	}

	if err = kubeHelper.CordonAndDrain(hostname); err != nil {
		return
	}

	return err
}

func upgradeRoot(url string) (err error) {
	// Identify next root disk
	var dev *probe.ProbedBlockDevice
	if dev, err = probe.GetDevWithFileSystemLabel(constants.NextRootPartitionLabel()); err != nil {
		return
	}

	// Should we handle anything around the disk/partition itself? maybe Format()
	rootTarget := install.Target{
		Label:          constants.NextRootPartitionLabel(),
		MountPoint:     "/var/nextroot",
		Device:         dev.Path,
		FileSystemType: "xfs",
		Force:          true,
		PartitionName:  dev.Path,
		Assets:         []*install.Asset{},
	}

	rootTarget.Assets = append(rootTarget.Assets, &install.Asset{
		Source:      url + "/" + constants.RootfsAsset,
		Destination: constants.RootfsAsset,
	})

	if err = rootTarget.Format(); err != nil {
		return
	}

	// mount up 'next' root disk rw
	mountpoint := mount.NewMountPoint(dev.Path, rootTarget.MountPoint, dev.SuperBlock.Type(), unix.MS_NOATIME, "")

	if err = mount.WithRetry(mountpoint); err != nil {
		return errors.Errorf("error mounting partitions: %v", err)
	}

	// Ensure we unmount the new rootfs in case of failure
	// nolint: errcheck
	defer mount.UnWithRetry(mountpoint)

	// install assets
	return rootTarget.Install()
}

func upgradeBoot(url string) error {
	bootTarget := install.Target{
		Label:      constants.BootPartitionLabel,
		MountPoint: "/boot",
		Assets:     []*install.Asset{},
	}

	// Kernel
	bootTarget.Assets = append(bootTarget.Assets, &install.Asset{
		Source:      url + "/" + constants.KernelAsset,
		Destination: filepath.Join(constants.NextRootPartitionLabel(), constants.KernelAsset),
	})

	// Initramfs
	bootTarget.Assets = append(bootTarget.Assets, &install.Asset{
		Source:      url + "/" + constants.InitramfsAsset,
		Destination: filepath.Join(constants.NextRootPartitionLabel(), constants.InitramfsAsset),
	})

	var err error
	if err = bootTarget.Install(); err != nil {
		return err
	}

	// TODO: Update gptmbr.bin

	// TODO: Figure out a method to update kernel args
	var cmdline []byte
	if cmdline, err = kernel.ReadProcCmdline(); err != nil {
		return err
	}

	// Create ExtlinuxConf struct
	extlinux := &syslinux.ExtlinuxConf{
		Default: constants.NextRootPartitionLabel(),
		Labels: []*syslinux.ExtlinuxConfLabel{
			{
				Root:   constants.NextRootPartitionLabel(),
				Kernel: filepath.Join(constants.NextRootPartitionLabel(), constants.KernelAsset),
				Initrd: filepath.Join(constants.NextRootPartitionLabel(), constants.InitramfsAsset),
				Append: string(cmdline),
			},
			// TODO see if we can omit this section and/or pass in only the root(?)
			// so we can make use of existing boot config
			{
				Root:   constants.CurrentRootPartitionLabel(),
				Kernel: filepath.Join(constants.CurrentRootPartitionLabel(), constants.KernelAsset),
				Initrd: filepath.Join(constants.CurrentRootPartitionLabel(), constants.InitramfsAsset),
				Append: string(cmdline),
			},
		},
	}

	return syslinux.Install(constants.BootMountPoint, extlinux)
}
