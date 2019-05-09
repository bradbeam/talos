/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

// nolint: dupl,golint
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/talos-systems/talos/cmd/osctl/pkg/client"
	"github.com/talos-systems/talos/cmd/osctl/pkg/helpers"
	"github.com/talos-systems/talos/internal/pkg/constants"
)

var (
	upgradeVersion string
	assetURL       string
)

// upgradeCmd represents the processes command
var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade Talos on the target node",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		creds, err := client.NewDefaultClientCredentials(talosconfig)
		if err != nil {
			helpers.Fatalf("error getting client credentials: %s", err)
		}
		if target != "" {
			creds.Target = target
		}
		c, err := client.NewClient(constants.OsdPort, creds)
		if err != nil {
			helpers.Fatalf("error constructing client: %s", err)
		}

		// TODO: See if we can validate version and prevent
		// starting upgrades to an unknown version
		if err := c.Upgrade(upgradeVersion, assetURL); err != nil {
			helpers.Fatalf("error upgrading host: %s", err)
		}
	},
}

func init() {
	upgradeCmd.Flags().StringVarP(&upgradeVersion, "version", "v", "", "target version to upgrade to")
	upgradeCmd.Flags().StringVarP(&assetURL, "url", "u", "", "url hosting upgrade assets (excluding filename)")
	upgradeCmd.Flags().StringVarP(&target, "target", "t", "", "target the specificed node")
	rootCmd.AddCommand(upgradeCmd)
}
