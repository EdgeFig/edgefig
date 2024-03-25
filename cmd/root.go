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
			Firewall: edgeconfig.Firewall{
				AllPing: edgeconfig.Enable,
			},
			Interfaces: map[string]edgeconfig.Interface{
				"eth0": {
					Enable: edgeconfig.Disable,
				},
				"eth1": {
					Enable: edgeconfig.Enable,
					Address: []netip.Prefix{
						netip.MustParsePrefix("10.0.0.3/22"),
						netip.MustParsePrefix("2001:db8::1/64"),
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
