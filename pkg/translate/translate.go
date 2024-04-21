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

	// Parse out firewall zones/rules
	for _, zoneYML := range router.Firewall.Zones {
		namePrefix := ""
		if zoneYML.IPType == types.IPAddressTypeV6 {
			namePrefix = "ipv6-"
		}
		_zone := edgeconfig.FirewallZone{
			NamePrefix:    namePrefix,
			Name:          zoneYML.Name,
			DefaultAction: zoneYML.DefaultAction,
			Description:   zoneYML.Description,
		}

		// Handles Rules
		for _, ruleYML := range zoneYML.Rules {
			_rule := edgeconfig.FirewallRule{
				Action:      ruleYML.Action,
				Description: ruleYML.Description,
				Destination: ruleYML.Destination,
				Log:         ruleYML.Log,
				Protocol:    ruleYML.Protocol,
				State: edgeconfig.FirewallRuleState{
					Established: ruleYML.Established,
					Invalid:     ruleYML.Invalid,
					New:         ruleYML.New,
					Related:     ruleYML.Related,
				},
			}

			_zone.Rules = append(_zone.Rules, _rule)
		}

		// Handles assignment of the zone to interfaces
		for _, zoneIfaceName := range zoneYML.In {
			for _ifaceIdx, _iface := range defaultRouter.Interfaces.Interfaces {
				if _iface.Name == zoneIfaceName {
					if zoneYML.IPType == types.IPAddressTypeV6 {
						if _iface.Firewall.In.V6Name != "" {
							panic("Duplicate zone assigned to interface")
						}
						defaultRouter.Interfaces.Interfaces[_ifaceIdx].Firewall.In.V6Name = zoneYML.Name
						break

					} else {
						if _iface.Firewall.In.Name != "" {
							panic("Duplicate zone assigned to interface")
						}
						defaultRouter.Interfaces.Interfaces[_ifaceIdx].Firewall.In.Name = zoneYML.Name
						break
					}
				}
			}
		}
		for _, zoneIfaceName := range zoneYML.Local {
			for _ifaceIdx, _iface := range defaultRouter.Interfaces.Interfaces {
				if _iface.Name == zoneIfaceName {
					if zoneYML.IPType == types.IPAddressTypeV6 {
						if _iface.Firewall.Local.V6Name != "" {
							panic("Duplicate zone assigned to interface")
						}
						defaultRouter.Interfaces.Interfaces[_ifaceIdx].Firewall.Local.V6Name = zoneYML.Name
						break
					} else {
						if _iface.Firewall.Local.Name != "" {
							panic("Duplicate zone assigned to interface")
						}
						defaultRouter.Interfaces.Interfaces[_ifaceIdx].Firewall.Local.Name = zoneYML.Name
						break
					}
				}
			}
		}
		for _, zoneIfaceName := range zoneYML.Out {
			for _ifaceIdx, _iface := range defaultRouter.Interfaces.Interfaces {
				if _iface.Name == zoneIfaceName {
					if zoneYML.IPType == types.IPAddressTypeV6 {
						if _iface.Firewall.Out.V6Name != "" {
							panic("Duplicate zone assigned to interface")
						}
						defaultRouter.Interfaces.Interfaces[_ifaceIdx].Firewall.Out.V6Name = zoneYML.Name
						break
					} else {
						if _iface.Firewall.Out.Name != "" {
							panic("Duplicate zone assigned to interface")
						}
						defaultRouter.Interfaces.Interfaces[_ifaceIdx].Firewall.Out.Name = zoneYML.Name
						break
					}
				}
			}
		}

		defaultRouter.Firewall.Zones = append(defaultRouter.Firewall.Zones, _zone)
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

		if staticRoute.NextHop.Is6() {
			edgeRouteConfig.RouteSuffix = "6"
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
		newUser := edgeconfig.User{
			Username: user.Username,
			Authentication: edgeconfig.Authentication{
				PlaintextPassword: user.Password,
			},
			Level: user.Role,
		}

		if user.Username == "ubnt" {
			// This is the user in defaults, so we need to overwrite item 0
			defaultRouter.System.Login.Users[0] = newUser
		} else {
			defaultRouter.System.Login.Users = append(defaultRouter.System.Login.Users, newUser)
		}

	}

	return defaultRouter, nil
}
