package translate

import (
	"fmt"

	"github.com/cmmarslender/edgefig/pkg/config"
	"github.com/cmmarslender/edgefig/pkg/edgeconfig"
)

// ConfigToEdgeConfig translates the friendly config to edgerouter config
// @TODO this should return a whole set of configs, not just router configs
func ConfigToEdgeConfig(cfg *config.Config) (*edgeconfig.Router, error) {
	edgecfg := &edgeconfig.Router{
		Firewall:   edgeconfig.Firewall{}, // @TODO
		Interfaces: edgeconfig.Interfaces{},
		Service:    edgeconfig.RouterServices{},
		System:     edgeconfig.RouterSystem{},
	}

	if len(cfg.Routers) == 0 {
		return nil, fmt.Errorf("no routers configured")
	}

	// @TODO Deal with more than the 0th router
	router := cfg.Routers[0]

	for intf, intCfg := range router.Interfaces {
		_iface := edgeconfig.Interface{
			Type:        edgeconfig.InterfaceTypeEthernet,
			Name:        intf,
			State:       edgeconfig.Enabled,
			Description: intCfg.Name,
			Address:     intCfg.Addresses,
			//Duplex:      "",
			//Speed:       "",
		}
		if intCfg.MTU != 0 {
			_iface.MTU = intCfg.MTU
		}

		edgecfg.Interfaces.Interfaces = append(edgecfg.Interfaces.Interfaces, _iface)
	}

	_dhcpServer := edgeconfig.DHCPServer{
		Disabled:       len(router.DHCP) == 0,
		HostfileUpdate: false,
		StaticARP:      false,
		UseDNSMASQ:     false,
	}
	for _, dhcpCfg := range router.DHCP {
		_dhcpNetwork := edgeconfig.DHCPNetwork{
			Name:          dhcpCfg.Name,
			Authoritative: true,
			Subnets: []edgeconfig.DHCPSubnet{
				{
					Subnet: dhcpCfg.Subnet,
					Router: dhcpCfg.Router,
					Lease:  dhcpCfg.Lease,
					DNS:    dhcpCfg.DNS,
					StartStop: edgeconfig.DHCPStartStop{
						Start: dhcpCfg.Start,
						Stop:  dhcpCfg.Stop,
					},
				},
			},
		}

		_dhcpServer.Networks = append(_dhcpServer.Networks, _dhcpNetwork)
	}
	edgecfg.Service.DHCPServer = _dhcpServer

	_natService := edgeconfig.NatService{Rules: []edgeconfig.NatRule{}}
	for _, natRule := range router.NAT {
		_natService.Rules = append(_natService.Rules, edgeconfig.NatRule{
			Name:              natRule.Name,
			Type:              natRule.Type,
			InboundInterface:  natRule.InboundInterface,
			OutboundInterface: natRule.OutboundInterface,
			Protocol:          natRule.Protocol,
			Log:               edgeconfig.EnableDisable(natRule.Log),
			OutsideAddress:    natRule.OutsideAddress,
			InsideAddress:     natRule.InsideAddress,
		})
	}
	edgecfg.Service.NAT = _natService

	edgecfg.System.HostName = router.Name
	for _, user := range router.Users {
		edgecfg.System.Login.Users = append(edgecfg.System.Login.Users, edgeconfig.User{
			Username: user.Username,
			Authentication: edgeconfig.Authentication{
				EncryptedPassword: user.Password,
			},
			Level: user.Role,
		})
	}

	return edgecfg, nil
}
