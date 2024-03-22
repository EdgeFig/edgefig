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
