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
					Type:  edgeconfig.InterfaceTypeEthernet,
					Name:  "eth0",
					State: edgeconfig.Enabled,
					Address: []netip.Prefix{
						netip.MustParsePrefix("192.168.1.1/24"),
					},
				},
				{
					Type:        edgeconfig.InterfaceTypeEthernet,
					Name:        "eth1",
					State:       edgeconfig.Enabled,
					AddressDHCP: "dhcp",
				},
				{
					Type:  edgeconfig.InterfaceTypeEthernet,
					Name:  "eth2",
					State: edgeconfig.Disabled,
				},
				{
					Type:  edgeconfig.InterfaceTypeEthernet,
					Name:  "eth3",
					State: edgeconfig.Disabled,
				},
				{
					Type:  edgeconfig.InterfaceTypeEthernet,
					Name:  "eth4",
					State: edgeconfig.Disabled,
				},
				{
					Type:  edgeconfig.InterfaceTypeEthernet,
					Name:  "eth5",
					State: edgeconfig.Disabled,
				},
				{
					Type:  edgeconfig.InterfaceTypeEthernet,
					Name:  "eth6",
					State: edgeconfig.Disabled,
				},
				{
					Type:  edgeconfig.InterfaceTypeEthernet,
					Name:  "eth7",
					State: edgeconfig.Disabled,
				},
				// @TODO loopback is a special case, and is an empty block
				// loopback lo {
				// }
				{
					Type:  edgeconfig.InterfaceTypeLoopback,
					Name:  "lo",
					State: edgeconfig.Enabled,
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
				OlderCiphers: edgeconfig.Disable, // Default is usually enabled
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
