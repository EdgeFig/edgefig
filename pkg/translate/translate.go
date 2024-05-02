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

		if len(intCfg.IPv6.Prefixes) > 0 {
			_iface.IPv6.DupAddrDetectTransmits = 1
			_iface.IPv6.RouterAdvert = edgeconfig.IPv6RouterAdvert{
				CurHopLimit:     64,
				LinkMTU:         intCfg.MTU,
				ManagedFlag:     false,
				MaxInterval:     600,
				NameServer:      intCfg.IPv6.Nameserver,
				OtherConfigFlag: false,
				ReachableTime:   0,
				RetransTimer:    0,
				SendAdvert:      true,
			}

			for _, prefixYml := range intCfg.IPv6.Prefixes {
				_iface.IPv6.RouterAdvert.Prefixes = append(_iface.IPv6.RouterAdvert.Prefixes, edgeconfig.IPv6PrefixAdvertisement{
					Prefix:         prefixYml.Prefix,
					AutonomousFlag: prefixYml.Autonomous,
					OnLinkFlag:     true,
					ValidLifetime:  2592000,
				})
			}
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

	// Parse out firewall groups
	for _, groupYML := range router.Firewall.Groups.AddressGroups {
		defaultRouter.Firewall.Group.AddressGroups = append(defaultRouter.Firewall.Group.AddressGroups, edgeconfig.AddressGroup{
			Name:        groupYML.Name,
			AddressPort: groupYML.AddressPort,
			Description: groupYML.Description,
		})
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
				Source: ruleYML.Source,
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
				RouterID: bgpCfg.RouterID.String(),
			},
		}

		for _, bgpPeer := range bgpCfg.Peers {
			bgpGroupNameTo := getBGPGroupName(bgpPeer, bgpDirTo)
			bgpGroupNameFrom := getBGPGroupName(bgpPeer, bgpDirFrom)

			prefixListFrom := edgeconfig.PrefixList{
				Name:  bgpGroupNameFrom,
				Rules: []edgeconfig.PrefixListRule{},
			}
			prefixListTo := edgeconfig.PrefixList{
				Name:  bgpGroupNameTo,
				Rules: []edgeconfig.PrefixListRule{},
			}
			routeMapFrom := edgeconfig.RouteMap{
				Name: bgpGroupNameFrom,
			}
			routeMapTo := edgeconfig.RouteMap{
				Name: bgpGroupNameTo,
			}

			if bgpPeer.IP.Is6() {
				prefixListFrom.PrefixListSuffix = "6"
				prefixListTo.PrefixListSuffix = "6"

				routeMapFrom.Rules = []edgeconfig.RouteMapRule{
					{
						Action: types.Permit,
						Match: edgeconfig.RouteMapMatch{
							IPv6: edgeconfig.RouteMatchIP{
								Address: edgeconfig.RouteMapAddress{
									PrefixList: bgpGroupNameFrom,
								},
							},
						},
					},
				}
				routeMapTo.Rules = []edgeconfig.RouteMapRule{
					{
						Action: types.Permit,
						Match: edgeconfig.RouteMapMatch{
							IPv6: edgeconfig.RouteMatchIP{
								Address: edgeconfig.RouteMapAddress{
									PrefixList: bgpGroupNameTo,
								},
							},
						},
					},
				}
			} else {
				routeMapFrom.Rules = []edgeconfig.RouteMapRule{
					{
						Action: types.Permit,
						Match: edgeconfig.RouteMapMatch{
							IPv4: edgeconfig.RouteMatchIP{
								Address: edgeconfig.RouteMapAddress{
									PrefixList: bgpGroupNameFrom,
								},
							},
						},
					},
				}
				routeMapTo.Rules = []edgeconfig.RouteMapRule{
					{
						Action: types.Permit,
						Match: edgeconfig.RouteMapMatch{
							IPv4: edgeconfig.RouteMatchIP{
								Address: edgeconfig.RouteMapAddress{
									PrefixList: bgpGroupNameTo,
								},
							},
						},
					},
				}
			}

			nbr := edgeconfig.BGPNeighbor{
				IP:  bgpPeer.IP,
				ASN: bgpPeer.ASN,
				DefaultOriginate: edgeconfig.BGPDefaultOriginate{
					Originate: bgpPeer.AnnounceDefault,
				},
				SoftReconfiguration: edgeconfig.BGPSoftReconfiguration{
					Inbound: types.KeyWhenEnabled(true),
				},
				UpdateSource: bgpPeer.SourceIP,
			}
			if bgpPeer.IP.Is6() {
				nbr.AddressFamily.IPv6Unicast.RouteMap = edgeconfig.BGPNeighborRouteMap{
					Export: bgpGroupNameTo,
					Import: bgpGroupNameFrom,
				}
			} else {
				nbr.RouteMap = edgeconfig.BGPNeighborRouteMap{
					Export: bgpGroupNameTo,
					Import: bgpGroupNameFrom,
				}
			}

			edgeBGPConfig.Neighbors = append(edgeBGPConfig.Neighbors, nbr)

			for _, prefix := range bgpPeer.Announcements {
				prefixListTo.Rules = append(prefixListTo.Rules, edgeconfig.PrefixListRule{
					Action: types.Permit,
					Prefix: prefix,
				})

				if prefix.Addr().Is6() {
					// Specifically add the network to the BGP section in the config
					edgeBGPConfig.AddressFamily.IPv6Unicast.Networks = append(edgeBGPConfig.AddressFamily.IPv6Unicast.Networks, edgeconfig.BGPNetwork{Prefix: prefix})
				} else {
					// Specifically add the network to the BGP section in the config
					edgeBGPConfig.Networks = append(edgeBGPConfig.Networks, edgeconfig.BGPNetwork{Prefix: prefix})
				}
			}

			for _, accept := range bgpPeer.Accept {
				prefixListFrom.Rules = append(prefixListFrom.Rules, edgeconfig.PrefixListRule{
					Action: types.Permit,
					GE:     accept.GE,
					LE:     accept.LE,
					Prefix: accept.Prefix,
				})
			}

			defaultRouter.Policy.PrefixLists = append(defaultRouter.Policy.PrefixLists, prefixListFrom, prefixListTo)
			defaultRouter.Policy.RouteMaps = append(defaultRouter.Policy.RouteMaps, routeMapFrom, routeMapTo)
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
				Interface:   staticRoute.Interface,
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
					Subnet:          dhcpCfg.Subnet,
					Router:          dhcpCfg.Router,
					Lease:           dhcpCfg.Lease,
					DNS:             dhcpCfg.DNS,
					Domain:          dhcpCfg.Domain,
					UnifiController: dhcpCfg.UnifiController,
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

	defaultRouter.Service.DNS.Forwarding.CacheSize = router.DNS.Forwarding.CacheSize
	defaultRouter.Service.DNS.Forwarding.ListenOn = router.DNS.Forwarding.ListenOn
	defaultRouter.Service.DNS.Forwarding.NameServers = router.DNS.Forwarding.Nameservers

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
