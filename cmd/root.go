package cmd

import (
	"fmt"
	"log"
	"net/netip"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cmmarslender/edgefig/pkg/edgeconfig"
	"github.com/cmmarslender/edgefig/pkg/types"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "edgefig",
	Short: "Configuration tool for edge* line of equipment from Ubiquiti",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		cfg := edgeconfig.Router{
			System: edgeconfig.RouterSystem{
				HostName: "EdgeRouter-Infinity",
				Login: edgeconfig.RouterLogin{
					Users: []edgeconfig.User{
						{
							Username: "ubnt",
							Authentication: edgeconfig.Authentication{
								// This is just the default "ubnt" password
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
			Firewall: edgeconfig.Firewall{
				AllPing:       edgeconfig.Enable,
				LogMartians:   edgeconfig.Disable,
				SendRedirects: edgeconfig.Enable,
				SynCookies:    edgeconfig.Enable,
			},
			Interfaces: edgeconfig.Interfaces{
				Interfaces: []edgeconfig.Interface{
					{
						Name:  "eth0",
						State: edgeconfig.Disabled,
					},
					{
						Name:        "eth1",
						State:       edgeconfig.Enabled,
						Description: "WAN",
						Address: []netip.Prefix{
							netip.MustParsePrefix("10.0.0.3/22"),
							netip.MustParsePrefix("2001:db8::1/64"),
						},
					},
				},
			},
			Service: edgeconfig.RouterServices{
				DHCPServer: edgeconfig.DHCPServer{
					Disabled:       false,
					HostfileUpdate: false,
					StaticARP:      false,
					UseDNSMASQ:     false,
					Networks: []edgeconfig.DHCPNetwork{
						{
							Name:          "LAN_10_100_0_X",
							Authoritative: edgeconfig.Enable,
							Subnets: []edgeconfig.DHCPSubnet{
								{
									Subnet: netip.MustParsePrefix("10.100.0.0/24"),
									Router: netip.MustParseAddr("10.100.0.1"),
									DNS: []netip.Addr{
										netip.MustParseAddr("1.1.1.1"),
										netip.MustParseAddr("8.8.8.8"),
									},
									Lease: 86400,
									StartStop: edgeconfig.DHCPStartStop{
										Start: netip.MustParseAddr("10.100.0.150"),
										Stop:  netip.MustParseAddr("10.100.0.254"),
									},
								},
							},
						},
					},
				},
				GUI: edgeconfig.GUIService{
					HTTPPort:     80,
					HTTPSPort:    443,
					OlderCiphers: edgeconfig.Disable,
				},
				NAT: edgeconfig.NatService{
					Dest: []edgeconfig.NatRule{
						{
							Name:             "Simple Port Forward",
							Type:             types.NATTypeDestination,
							InboundInterface: "eth1",
							Protocol:         types.ProtocolTCP,
							Log:              edgeconfig.Disable,
							OutsideAddress: types.NATAddress{
								Address: netip.MustParseAddr("192.168.1.1"),
								Port:    443,
							},
							InsideAddress: types.NATAddress{
								Address: netip.MustParseAddr("10.48.0.50"),
								Port:    443,
							},
						},
						{
							Name:             "Simple IP Forward",
							Type:             types.NATTypeDestination,
							InboundInterface: "eth1",
							Log:              edgeconfig.Disable,
							OutsideAddress: types.NATAddress{
								Address: netip.MustParseAddr("192.168.1.2"),
							},
							InsideAddress: types.NATAddress{
								Address: netip.MustParseAddr("10.48.0.51"),
							},
						},
					},
					Src: []edgeconfig.NatRule{
						{
							Name:              "Masquerade for WAN",
							Log:               edgeconfig.Disable,
							OutboundInterface: "eth1",
							Protocol:          types.ProtocolAll,
							Type:              types.NATTypeMasquerade,
						},
					},
				},
				SSH: edgeconfig.SSHService{
					Port:            22,
					ProtocolVersion: "v2",
				},
			},
		}

		data, err := edgeconfig.Marshal(cfg)
		if err != nil {
			log.Fatalln(err.Error())
		}

		fmt.Printf("%s", string(data))
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.edgefig.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".edgefig" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".edgefig")
	}

	viper.SetEnvPrefix("EDGEFIG_")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
