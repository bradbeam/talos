/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package reg

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/talos-systems/talos/internal/app/osd/proto"
	"github.com/talos-systems/talos/internal/pkg/upgrade"
)

// Upgrade initiates a Talos upgrade ... and implements the proto.OSDServer
// interface
func (r *Registrator) Upgrade(ctx context.Context, in *proto.UpgradeRequest) (data *proto.UpgradeReply, err error) {

	if err = upgrade.NewUpgrade(in.Url); err != nil {
		return data, err
	}

	// Trigger reboot
	defer func() {
		if _, err = r.Reboot(ctx, &empty.Empty{}); err != nil {
			return
		}
	}()

	// profit
	data = &proto.UpgradeReply{Ack: "Upgrade initiated"}
	return data, err
}
