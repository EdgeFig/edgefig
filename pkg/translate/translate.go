package translate

import (
	"fmt"

	"github.com/cmmarslender/edgefig/pkg/config"
	"github.com/cmmarslender/edgefig/pkg/edgeconfig"
	"github.com/cmmarslender/edgefig/pkg/types"
)

// ConfigToEdgeConfig translates the friendly config to edgerouter config
// @TODO this should return a whole set of configs, not just router configs
func ConfigToEdgeConfig(cfg *config.Config) (*edgeconfig.Router, error) {
	if len(cfg.Routers) == 0 {
		return nil, fmt.Errorf("no routers configured")
	}

	// @TODO Deal with more than the 0th router
	router := cfg.Routers[0]

	defaultRouter := getDefaultRouterConfig()
	defaultRouter.Firewall.AllPing = types.Enable
	defaultRouter.Firewall.SendRedirects = types.Enable
	defaultRouter.Firewall.SynCookies = types.Enable

	// Parse out firewall zones/rules
	for zoneName, zoneYML := range router.Firewall.Zones {
		namePrefix := ""
		if zoneYML.IPType == types.IPAddressTypeV6 {
			namePrefix = "ipv6-"
		}
		_zone := edgeconfig.FirewallZone{
			NamePrefix:    namePrefix,
			Name:          zoneName,
			DefaultAction: zoneYML.DefaultAction,
			Description:   zoneYML.Description,
			Rules:         []edgeconfig.FirewallRule{},
		}

		// @TODO handle rules

		defaultRouter.Firewall.Zones = append(defaultRouter.Firewall.Zones, _zone)
	}

	for intf, intCfg := range router.Interfaces {
		_iface := edgeconfig.Interface{
			Name:        intf,
			State:       types.Enabled,
			Description: intCfg.Name,
			Address:     intCfg.Addresses,
			Speed:       edgeconfig.AutoUint32(intCfg.Speed),
			Duplex:      edgeconfig.AutoString(intCfg.Duplex),
		}
		if intCfg.MTU != 0 {
			_iface.MTU = intCfg.MTU
		}

		for _, vlanName := range intCfg.VLANs {
			vlanCfg, err := cfg.GetVLANByName(vlanName)
			if err != nil {
				return nil, err
			}

			edgeVlan := edgeconfig.VLAN{
				ID:          vlanCfg.ID,
				Address:     vlanCfg.Address,
				Description: vlanCfg.Name,
			}
			if vlanCfg.MTU > 0 {
				edgeVlan.MTU = vlanCfg.MTU
			}
			_iface.VLANs = append(_iface.VLANs, edgeVlan)
		}

		// @TODO make some methods to keep references by key vs this hunting/replacing
		// This iterates the default interfaces and injects our customized config
		// Since we have to have all interfaces defined, this was an easy way to accomplish that
		for replI, replInt := range defaultRouter.Interfaces.Interfaces {
			if replInt.Name == _iface.Name {
				defaultRouter.Interfaces.Interfaces[replI] = _iface
			}
		}
	}

	for _, bgpCfg := range router.BGP {
		edgeBGPConfig := edgeconfig.BGPConfig{
			ASN:       bgpCfg.ASN,
			Neighbors: make([]edgeconfig.BGPNeighbor, 0),
			Parameters: edgeconfig.BGPParameters{
				RouterID: bgpCfg.IP.String(),
			},
		}

		for _, bgpPeer := range bgpCfg.Peers {
			nbr := edgeconfig.BGPNeighbor{
				IP:  bgpPeer.IP,
				ASN: bgpPeer.ASN,
				DefaultOriginate: edgeconfig.BGPDefaultOriginate{
					Originate: bgpPeer.AnnounceDefault,
				},
				SoftReconfiguration: edgeconfig.BGPSoftReconfiguration{
					Inbound: types.KeyWhenEnabled(true),
				},
			}

			edgeBGPConfig.Neighbors = append(edgeBGPConfig.Neighbors, nbr)
		}

		for _, prefix := range bgpCfg.Announcements {
			edgeBGPConfig.Networks = append(edgeBGPConfig.Networks, edgeconfig.BGPNetwork{Prefix: prefix})
		}

		defaultRouter.Protocols.BGP = append(defaultRouter.Protocols.BGP, edgeBGPConfig)
	}

	for _, staticRoute := range router.Routes {
		edgeRouteConfig := edgeconfig.StaticRoute{
			Route: staticRoute.Route,
			NextHop: edgeconfig.NextHop{
				NextHop:     staticRoute.NextHop,
				Description: staticRoute.Description,
				Distance:    staticRoute.Distance,
			},
		}

		defaultRouter.Protocols.Static.Routes = append(defaultRouter.Protocols.Static.Routes, edgeRouteConfig)
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

		for _, reservation := range dhcpCfg.Reservations {
			_dhcpNetwork.Subnets[0].StaticMappings = append(
				_dhcpNetwork.Subnets[0].StaticMappings,
				edgeconfig.DHCPStaticMapping{
					Name:       reservation.Name,
					IPAddress:  reservation.IP,
					MACAddress: reservation.MAC,
				},
			)
		}

		_dhcpServer.Networks = append(_dhcpServer.Networks, _dhcpNetwork)
	}
	defaultRouter.Service.DHCPServer = _dhcpServer

	_natService := edgeconfig.NatService{}
	for _, natRule := range router.NAT {
		newRule := edgeconfig.NatRule{
			Name:              natRule.Name,
			Type:              natRule.Type,
			InboundInterface:  natRule.InboundInterface,
			OutboundInterface: natRule.OutboundInterface,
			Protocol:          natRule.Protocol,
			Log:               types.EnableDisable(natRule.Log),
		}
		if newRule.Type == types.NATTypeDestination {
			newRule.Destination = natRule.OutsideAddress
			newRule.InsideAddress = natRule.InsideAddress
			_natService.Dest = append(_natService.Dest, newRule)
		} else {
			newRule.Source = natRule.InsideAddress
			newRule.OutsideAddress = natRule.OutsideAddress
			_natService.Src = append(_natService.Src, newRule)
		}

	}
	defaultRouter.Service.NAT = _natService

	defaultRouter.System.HostName = router.Name
	for _, user := range router.Users {
		defaultRouter.System.Login.Users = append(defaultRouter.System.Login.Users, edgeconfig.User{
			Username: user.Username,
			Authentication: edgeconfig.Authentication{
				EncryptedPassword: user.Password,
			},
			Level: user.Role,
		})
	}

	return defaultRouter, nil
}
