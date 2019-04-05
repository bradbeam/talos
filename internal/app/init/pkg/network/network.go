/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package network

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/autonomy/dhcp/dhcpv4"
	"github.com/autonomy/dhcp/dhcpv4/client4"
	"github.com/autonomy/dhcp/netboot"
	"github.com/autonomy/talos/internal/pkg/userdata"
	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
)

// Setup creates the network.
func Setup(data *userdata.UserData) (err error) {

	// If no networking config is defined,
	// bring up lo and eth0 with dhcp on eth0
	if data == nil || data.Networking.OS == nil {
		log.Println("default network setup")
		return defaultNetworkSetup()
	}

	// TODO: Turn this into a log level
	/*
		log.Println("All available network links")
		links, _ := netlink.LinkList()
		for _, link := range links {
			log.Printf("%+v", link)
		}
	*/

	// Always bring up lo by default
	log.Println("bringing up lo")
	if err = ifup("lo"); err != nil {
		return err
	}

	// Iterate through defined network devices
	log.Println("starting up network devices")
	for _, netconf := range data.Networking.OS.Devices {
		// Normal Interface
		if netconf.Bond == nil {
			log.Println("bringing up normal interface")
			if err := ifup(netconf.Interface); err != nil {
				return err
			}
		} else {
			// TODO test
			log.Println("bringing up bonded interface")
			bond := netlink.NewLinkBond(netlink.LinkAttrs{Name: netconf.Interface})
			if _, ok := netlink.StringToBondModeMap[netconf.Bond.Mode]; !ok {
				return errors.New(fmt.Sprintf("invalid bond mode for s", netconf.Interface))
			}
			bond.Mode = netlink.StringToBondModeMap[netconf.Bond.Mode]

			if _, ok := netlink.StringToBondLacpRateMap[netconf.Bond.LACPRate]; !ok {
				return errors.New(fmt.Sprintf("invalid lacp rate for %s", netconf.Interface))
			}
			bond.LacpRate = netlink.StringToBondLacpRateMap[netconf.Bond.LACPRate]

			if _, ok := netlink.StringToBondXmitHashPolicyMap[netconf.Bond.HashPolicy]; !ok {
				return errors.New(fmt.Sprintf("invalid lacp rate for %s", netconf.Interface))
			}
			bond.XmitHashPolicy = netlink.StringToBondXmitHashPolicyMap[netconf.Bond.HashPolicy]

			// Set up bonding if defined
			var slaveLink netlink.Link
			for _, bondInterface := range netconf.Bond.Interfaces {
				log.Printf("enslaving %s for %s\n", bondInterface, netconf.Interface)
				slaveLink, err = netlink.LinkByName(bondInterface)
				if err != nil {
					return err
				}

				if err := netlink.LinkSetBondSlave(slaveLink, &netlink.Bond{LinkAttrs: *bond.Attrs()}); err != nil {
					return err
				}
			}
		}

		if netconf.DHCP {
			// Set up dhcp renewals every 5m
			go func() {
				for {
					// TODO pick this out of the dhclient/netconf response
					// so we can request less frequently
					// Example response:  #ValidLifetime
					// [   12.399713] [talos] [initramfs] netconf: &{Addresses:[{IPNet:{IP:147.75.64.19 Mask:fffffffe} PreferredLifetime:0 ValidLifetime:1800}] Classless:[] DNSServers:[147.75.207.207 147.75.207.208] DNSSearchList:[] Routers:[147.75.64.18]}
					time.Sleep(5 * time.Minute)
					log.Println("renewing dhcp lease")
					if err = dhclient(netconf.Interface); err != nil {
						// Probably need to do something better here but not sure there's much to do
						log.Println("failed to renew dhcp lease, ", err)
					}
				}
			}()
			// single/initial dhcp request
			if err = dhclient(netconf.Interface); err != nil {
				return err
			}
		} else {
			addr, _ := netlink.ParseAddr(netconf.CIDR)
			var link netlink.Link
			if link, err = netlink.LinkByName(netconf.Interface); err != nil {
				return err
			}
			if err = netlink.AddrAdd(link, addr); err != nil {
				return err
			}
		}
	}

	return nil
}

func dhclient(ifname string) error {
	var err error
	var netconf *netboot.NetConf

	// TODO: Figure out how we want to pass around ntp servers
	modifiers := []dhcpv4.Modifier{
		dhcpv4.WithRequestedOptions(
			dhcpv4.OptionHostName,
			dhcpv4.OptionClasslessStaticRouteOption,
			dhcpv4.OptionDNSDomainSearchList,
			dhcpv4.OptionNTPServers,
		),
	}

	if netconf, err = dhclient4(ifname, modifiers...); err != nil {
		return err
	}
	if err = netboot.ConfigureInterface(ifname, netconf); err != nil {
		return err
	}

	return err
}

// nolint: gocyclo
func dhclient4(ifname string, modifiers ...dhcpv4.Modifier) (*netboot.NetConf, error) {
	attempts := 10
	client := client4.NewClient()
	var (
		conv []*dhcpv4.DHCPv4
		err  error
	)
	for attempt := 0; attempt < attempts; attempt++ {
		log.Printf("requesting DHCP lease: attempt %d of %d", attempt+1, attempts)
		conv, err = client.Exchange(ifname, modifiers...)
		if err != nil && attempt < attempts {
			log.Printf("failed to request DHCP lease: %v", err)
			time.Sleep(time.Duration(attempt) * time.Second)
			continue
		}
		break
	}

	for _, m := range conv {
		if m.OpCode == dhcpv4.OpcodeBootReply && m.MessageType() == dhcpv4.MessageTypeOffer {
			if m.YourIPAddr != nil {
				log.Printf("using IP address %s", m.YourIPAddr.String())
			}

			hostname := m.YourIPAddr.String()
			if m.HostName() != "" {
				hostname = m.HostName()
			}
			log.Printf("using hostname: %s", hostname)
			if err = unix.Sethostname([]byte(hostname)); err != nil {
				return nil, err
			}

			break
		}
	}

	netconf, _, err := netboot.ConversationToNetconfv4(conv)
	if err != nil {
		return nil, err
	}

	return netconf, err
}

func ifup(ifname string) (err error) {
	var link netlink.Link
	if link, err = netlink.LinkByName(ifname); err != nil {
		return err
	}
	if err = netlink.LinkSetUp(link); err != nil {
		return err
	}
	return nil
}

func defaultNetworkSetup() (err error) {
	if err := ifup("lo"); err != nil {
		return err
	}
	if err = ifup("eth0"); err != nil {
		return err
	}

	return dhclient("eth0")
}
