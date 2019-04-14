package network

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/insomniacslk/dhcp/dhcpv4/nclient4"
	"github.com/jsimonetti/rtnetlink"
	"github.com/talos-systems/dhcp/dhcpv4"
)

const defaultInterface = eth0

type NetworkInterface struct {
	Name string
	Link *rtnetlink.LinkMessage
}

func Default(ifname string) (err error) {

	nic := NetworkInterface{Name: defaultInterface}

	if ifname != "" {
		nic.Name = ifname
	}

	rt := &RTNL{}
	defer rt.Close()

	// List out the network interfaces
	var links []rtnetlink.LinkMessage
	links, err = rt.List()
	if err != nil {
		return
	}

	// Verify the interface exists and is an appropriate type
	for _, iface := range links {
		if iface.Attributes.Name == nic.Name {
			switch iface.Type {
			case unix.ARPHRD_LOOPBACK:
			case unix.ARPHRD_ETHER:
			default:
				return errors.New(fmt.Sprintf("unsupported interface family %s", iface.Type))
			}
			nic.Link = iface
			break
		}
		return errors.New(fmt.Sprintf("interface %s not found", nic.Name))
	}

	// Bring up the interface if it's down
	// Do we need to add || rtnetlink.OperStateUnknown as well?
	if nic.Link.Attributes.OperationalState == rtnetlink.OperStateDown {
		if err = rt.LinkSet(int(nic.Index), true); err != nil {
			return
		}
	}

	// Get DHCP lease for primary interface
	var dhclient *nclient4.Client
	dhclient, err = nclient4.New(nic.Name)
	if err != nil {
		return err
	}

	modifiers := []dhcpv4.Modifier{
		dhcpv4.WithRequestedOptions(
			dhcpv4.OptionHostName,
			dhcpv4.OptionClasslessStaticRouteOption,
			dhcpv4.OptionDNSDomainSearchList,
			dhcpv4.OptionNTPServers,
		),
	}

	var offer, ack *dhcpv4.DHCPv4
	offer, ack, err = dhclient.Request(context.TODO, modifiers)
	if err != nil {
		return err
	}

	log.Printf("Offer: %+v", offer)
	log.Printf("ACK: %+v", ack)

	return
}

func (r *RTNL) Link() {
	// Move link list here
	// Move loop/link list validation here
}

func Setup(data *userdata.Data) (err error) {

}
