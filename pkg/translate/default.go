package translate

import (
	"fmt"
	"net/netip"
	"strings"

	"github.com/cmmarslender/edgefig/pkg/edgeconfig"
	"github.com/cmmarslender/edgefig/pkg/types"
)

func getDefaultRouterConfig(interfaces map[string]struct{}) *edgeconfig.Router {
	cfg := edgeconfig.Router{
		Interfaces: edgeconfig.Interfaces{
			// Edgerouters have eth0 on a static IP and eth1 on DHCP out of the box
			// We do discovery to determine what other interfaces exist
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
			},
			Switches: []edgeconfig.SwitchInterface{},
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

	skip := map[string]struct{}{"eth0": {}, "eth1": {}}
	for iface := range interfaces {
		if _, ok := skip[iface]; ok {
			continue
		}
		if strings.Contains(iface, "switch") {
			cfg.Interfaces.Switches = append(cfg.Interfaces.Switches, edgeconfig.SwitchInterface{
				Name: iface,
				MTU:  1500,
			})
		} else {
			cfg.Interfaces.Interfaces = append(cfg.Interfaces.Interfaces, edgeconfig.Interface{
				Name:  iface,
				State: types.Disabled,
			})
		}

	}

	err := cfg.Validate()
	if err != nil {
		panic(fmt.Sprintf("Defaults do not validate! %s", err.Error()))
	}

	return &cfg
}
