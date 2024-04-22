package translate

import (
	"fmt"
	"net/netip"

	"github.com/cmmarslender/edgefig/pkg/edgeconfig"
	"github.com/cmmarslender/edgefig/pkg/types"
)

func getDefaultRouterConfig() *edgeconfig.Router {
	cfg := edgeconfig.Router{
		Interfaces: edgeconfig.Interfaces{
			// @TODO need to do interface discovery to see how many interfaces are supported
			Interfaces: []edgeconfig.Interface{
				{
					Name:  "eth0",
					State: types.Enabled,
					Address: []netip.Prefix{
						netip.MustParsePrefix("192.168.1.1/24"),
					},
				},
				{
					Name:        "eth1",
					State:       types.Enabled,
					AddressDHCP: "dhcp",
				},
				{
					Name:  "eth2",
					State: types.Disabled,
				},
				{
					Name:  "eth3",
					State: types.Disabled,
				},
				{
					Name:  "eth4",
					State: types.Disabled,
				},
				{
					Name:  "eth5",
					State: types.Disabled,
				},
				{
					Name:  "eth6",
					State: types.Disabled,
				},
				{
					Name:  "eth7",
					State: types.Disabled,
				},
				{
					Name:  "eth8",
					State: types.Disabled,
				},
			},
		},
		Policy: edgeconfig.RouterPolicy{
			PrefixLists: []edgeconfig.PrefixList{
				{
					Name: RouteMapV4From,
					Rules: []edgeconfig.PrefixListRule{
						{
							Action: types.Permit,
							LE:     24,
							Prefix: netip.MustParsePrefix("0.0.0.0/0"),
						},
					},
				},
				{
					Name: RouteMapV4To,
				},
				{
					PrefixListSuffix: "6",
					Name:             RouteMapV6From,
					Rules: []edgeconfig.PrefixListRule{
						{
							Action: types.Permit,
							LE:     64,
							Prefix: netip.MustParsePrefix("0::/0"),
						},
					},
				},
				{
					PrefixListSuffix: "6",
					Name:             RouteMapV6To,
				},
			},
			RouteMaps: []edgeconfig.RouteMap{
				{
					Name: RouteMapV4From,
					Rules: []edgeconfig.RouteMapRule{
						{
							Action: types.Permit,
							Match: edgeconfig.RouteMapMatch{
								IPv4: edgeconfig.RouteMatchIP{
									Address: edgeconfig.RouteMapAddress{
										PrefixList: RouteMapV4From,
									},
								},
							},
						},
					},
				},
				{
					Name: RouteMapV4To,
					Rules: []edgeconfig.RouteMapRule{
						{
							Action: types.Permit,
							Match: edgeconfig.RouteMapMatch{
								IPv4: edgeconfig.RouteMatchIP{
									Address: edgeconfig.RouteMapAddress{
										PrefixList: RouteMapV4To,
									},
								},
							},
						},
					},
				},
				{
					Name: RouteMapV6From,
					Rules: []edgeconfig.RouteMapRule{
						{
							Action: types.Permit,
							Match: edgeconfig.RouteMapMatch{
								IPv6: edgeconfig.RouteMatchIP{
									Address: edgeconfig.RouteMapAddress{
										PrefixList: RouteMapV6From,
									},
								},
							},
						},
					},
				},
				{
					Name: RouteMapV6To,
					Rules: []edgeconfig.RouteMapRule{
						{
							Action: types.Permit,
							Match: edgeconfig.RouteMapMatch{
								IPv6: edgeconfig.RouteMatchIP{
									Address: edgeconfig.RouteMapAddress{
										PrefixList: RouteMapV6To,
									},
								},
							},
						},
					},
				},
			},
		},
		System: edgeconfig.RouterSystem{
			HostName: "EdgeRouter-Infinity", // @TODO this should be based on the detected router model
			Login: edgeconfig.RouterLogin{
				Users: []edgeconfig.User{
					{
						Username: "ubnt",
						Authentication: edgeconfig.Authentication{
							EncryptedPassword: "$1$zKNoUbAo$gomzUbYvgyUMcD436Wo66.",
						},
						Level: types.UserLevelAdmin,
					},
				},
			},
			NTP: edgeconfig.NTPServers{
				Servers: []edgeconfig.NTPServer{
					{
						Hostname: "0.ubnt.pool.ntp.org",
					},
					{
						Hostname: "1.ubnt.pool.ntp.org",
					},
					{
						Hostname: "2.ubnt.pool.ntp.org",
					},
					{
						Hostname: "3.ubnt.pool.ntp.org",
					},
				},
			},
			Syslog: edgeconfig.Syslog{
				Global: edgeconfig.SyslogGlobal{
					Facilities: []edgeconfig.SyslogFacility{
						{
							Name:  "all",
							Level: "notice",
						},
						{
							Name:  "protocols",
							Level: "debug",
						},
					},
				},
			},
			TimeZone: "UTC",
		},
		Service: edgeconfig.RouterServices{
			GUI: edgeconfig.GUIService{
				HTTPPort:     80,
				HTTPSPort:    443,
				OlderCiphers: types.Disable, // Default is usually enabled
			},
			SSH: edgeconfig.SSHService{
				Port:            22,
				ProtocolVersion: "v2",
			},
		},
	}

	err := cfg.Validate()
	if err != nil {
		panic(fmt.Sprintf("Defaults do not validate! %s", err.Error()))
	}

	return &cfg
}
