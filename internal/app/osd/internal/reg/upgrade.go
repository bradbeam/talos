/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package reg

import (
	"context"
	"os"
	"path"
	"path/filepath"

	"github.com/coreos/go-semver/semver"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	"github.com/talos-systems/talos/internal/app/osd/proto"
	"github.com/talos-systems/talos/internal/pkg/blockdevice/bootloader/syslinux"
	"github.com/talos-systems/talos/internal/pkg/blockdevice/probe"
	"github.com/talos-systems/talos/internal/pkg/constants"
	"github.com/talos-systems/talos/internal/pkg/install"
	"github.com/talos-systems/talos/internal/pkg/kernel"
	"github.com/talos-systems/talos/internal/pkg/kubernetes"
	"github.com/talos-systems/talos/internal/pkg/mount"
	"github.com/talos-systems/talos/internal/pkg/version"
	"golang.org/x/sys/unix"
)

// Upgrade initiates a Talos upgrade ... and implements the proto.OSDServer
// interface
func (r *Registrator) Upgrade(ctx context.Context, in *proto.UpgradeRequest) (data *proto.UpgradeReply, err error) {
	// Compare current running talos version against specified version
	newVersion := semver.New(in.Version)
	currentVersion := semver.New(version.Tag)

	// -1 less than | 0 equal | 1 greater than
	if newVersion.Compare(*currentVersion) != 1 {
		return
	}

	if err = upgradeRoot(in); err != nil {
		return
	}
	if err = upgradeBoot(in); err != nil {
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

	// reboot
	defer func() {
		if _, err = r.Reboot(ctx, &empty.Empty{}); err != nil {
			return
		}
	}()

	// profit
	data = &proto.UpgradeReply{Ack: "Upgrade initiated"}
	return data, err
}

func upgradeRoot(in *proto.UpgradeRequest) (err error) {
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
		PartitionName:  dev.Path,
		Assets:         []*install.Asset{},
	}

	rootTarget.Assets = append(rootTarget.Assets, &install.Asset{
		Source:      path.Join(in.Url, constants.RootfsAsset),
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

	// install assets
	return rootTarget.Install()
}

func upgradeBoot(in *proto.UpgradeRequest) error {
	bootTarget := install.Target{
		Label:      constants.BootPartitionLabel,
		MountPoint: "/boot",
		Assets:     []*install.Asset{},
	}

	// Kernel
	bootTarget.Assets = append(bootTarget.Assets, &install.Asset{
		Source:      path.Join(in.Url, constants.KernelAsset),
		Destination: filepath.Join(constants.NextRootPartitionLabel(), constants.KernelAsset),
	})

	// Initramfs
	bootTarget.Assets = append(bootTarget.Assets, &install.Asset{
		Source:      path.Join(in.Url, constants.InitramfsAsset),
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
