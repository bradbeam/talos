/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package networkd

import (
	"context"
	"log"
	"net"
	"sync"
	"time"

	"github.com/jsimonetti/rtnetlink"
	"github.com/jsimonetti/rtnetlink/rtnl"
	"github.com/talos-systems/talos/internal/app/networkd/pkg/address"
	"github.com/talos-systems/talos/internal/app/networkd/pkg/nic"
)

// Set up default nameservers
const (
	DefaultPrimaryResolver   = "1.1.1.1"
	DefaultSecondaryResolver = "8.8.8.8"
)

// Networkd provides the high level interaction to configure network interfaces
// on a host system. This currently support addressing configuration via dhcp
// and/or a specified configuration file.
type Networkd struct {
	Conn   *rtnl.Conn
	nlConn *rtnetlink.Conn
}

// New instantiates a new rtnetlink connection that is used for all subsequent
// actions
func New() (*Networkd, error) {
	// Handle netlink connection
	conn, err := rtnl.Dial(nil)
	if err != nil {
		return nil, err
	}

	// Need rtnetlink for MTU setting
	// TODO: possible rtnl enhancement
	nlConn, err := rtnetlink.Dial(nil)
	if err != nil {
		return nil, err
	}

	return &Networkd{Conn: conn, nlConn: nlConn}, err
}

// Discover enumerates a list of network links on the host and creates a
// base set of interface configuration options
func (n *Networkd) Discover() (NetConf, error) {
	links, err := n.Conn.Links()
	if err != nil {
		return NetConf{}, err
	}

	linkmap := NetConf{}

	for _, link := range filterInterfaceByName(links) {
		linkmap[link] = parseLinkMessage(link)
	}

	return linkmap, nil
}

// Configure handles the interface configuration portion. This is inclusive of
// the address discovery ( static vs dhcp ) as well as the netlink interaction
// to set an address on the link and create any routes.
func (n *Networkd) Configure(ifaces ...*nic.NetworkInterface) error {
	var (
		err       error
		resolvers []net.IP
		link      *net.Interface
	)

	for _, iface := range ifaces {
		log.Printf("configuring %+v\n", iface)

		// Attempt dhcp against all unconfigured interfaces
		if len(iface.AddressMethod) == 0 {
			link, err = n.Conn.LinkByIndex(int(iface.Index))
			if err != nil {
				log.Printf("failed to find interface %d", iface.Index)
				continue
			}
			iface.AddressMethod = append(iface.AddressMethod, &address.DHCP{NetIf: link})
		}

		// Bring up the interface
		if err = n.Conn.LinkUp(&net.Interface{Index: int(iface.Index)}); err != nil {
			log.Printf("failed to bring up %s: %v", iface.Name, err)
			continue
		}

		// Generate rtnetlink.AddressMessage for each address method defined on
		// the interface
		for _, method := range iface.AddressMethod {
			log.Printf("configuring %s addressing for %s\n", method.Name(), iface.Name)
			if err = n.configureInterface(method); err != nil {
				// Treat as non fatal error when failing to configure an interface
				log.Println(err)
				continue
			}

			// Aggregate a list of DNS servers/resolvers
			resolvers = append(resolvers, method.Resolvers()...)
		}
	}

	// Write out resolv.conf
	if err = writeResolvConf(resolvers); err != nil {
		return err
	}

	return nil
}

// Renew sets up a long running loop to refresh a network interfaces
// addressing configuration. Currently this only applies to interfaces
// configured by DHCP.
func (n *Networkd) Renew(ifaces ...*nic.NetworkInterface) {
	var wg sync.WaitGroup
	for _, iface := range ifaces {
		for _, method := range iface.AddressMethod {
			if method.TTL() == 0 {
				continue
			}
			wg.Add(1)

			go n.renew(method)
		}
	}

	// TODO switch this out with context when that gets added
	select {}
}

// renew sets up the looping to ensure we keep the addressing information
// up to date. We attempt to do our first reconfiguration halfway through
// address TTL. If that fails, we'll continue to attempt to retry every
// halflife.
func (n *Networkd) renew(method address.Addressing) {
	renewDuration := method.TTL() / 2
	for {
		<-time.After(renewDuration)

		if err := n.configureInterface(method); err != nil {
			log.Printf("failed to renew interface address for %s: %v\n", method.Link().Name, err)
			renewDuration = (renewDuration / 2)
		} else {
			renewDuration = method.TTL() / 2
		}
	}
}

// configureInterface handles the actual address discovery mechanism and
// netlink interaction to configure the interface
func (n *Networkd) configureInterface(method address.Addressing) error {
	// TODO s/Discover/Something else/
	// TODO make context more relevant
	var err error
	if err = method.Discover(context.Background()); err != nil {
		// Right now this would only happen during dhcp discovery failure
		log.Printf("failed to prep %s: %v", method.Link().Name, err)
		return err
	}

	// Set link MTU if we got a response
	if err = n.setMTU(method.Link().Index, method.MTU()); err != nil {
		log.Printf("failed to set mtu %d for %s: %v", method.MTU(), method.Link().Name, err)
		return err
	}

	// Check to see if we need to configure the address
	addrs, err := n.Conn.Addrs(method.Link(), method.Family())
	if err != nil {
		return err
	}

	addressExists := false
	for _, addr := range addrs {
		if method.Address().String() == addr.String() {
			addressExists = true
		}
	}

	if !addressExists {
		if err = n.Conn.AddrAdd(method.Link(), method.Address()); err != nil {
			log.Printf("failed to add address %+v to %s: %v", method.Address(), method.Link().Name, err)
			return err
		}
	}

	// Add any routes
	for _, r := range method.Routes() {
		if err = n.Conn.RouteAdd(method.Link(), *r.Dest, r.Router); err != nil {
			// TODO may need to check error type for EEXIST and skip
			// TODO how do we want to handle failures in routing?
			log.Printf("Failed to add route %+v for %s: %v", r, method.Link().Name, err)
			continue
		}
	}

	return err
}

// Hostname returns the first hostname found from the addressing methods.
func (n *Networkd) Hostname(ifaces ...*nic.NetworkInterface) string {
	for _, iface := range ifaces {
		for _, method := range iface.AddressMethod {
			if method.Hostname() != "" {
				return method.Hostname()
			}
		}
	}

	return ""
}

/*
// TODO add this in with some debug level of loggin
func (n *Networkd) printState() {
	rl, err := n.Conn.Route.List()
	if err != nil {
		log.Println(err)
		return
	}
	for _, r := range rl {
		log.Printf("%+v", r)
	}

	links, err := n.Conn.Link.List()
	if err != nil {
		log.Println(err)
		return
	}
	for _, link := range links {
		log.Printf("%+v", link)
		log.Printf("%+v", link.Attributes)
	}

	b, err := ioutil.ReadFile("/etc/resolv.conf")
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("resolv.conf: %s", string(b))
}
*/
