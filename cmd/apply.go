package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cmmarslender/edgefig/internal/connection"
	"github.com/cmmarslender/edgefig/internal/util"
	"github.com/cmmarslender/edgefig/pkg/config"
	"github.com/cmmarslender/edgefig/pkg/edgeconfig"
	"github.com/cmmarslender/edgefig/pkg/translate"
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Applies the configuration to all devices",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig(viper.GetString("config"))
		if err != nil {
			log.Fatalln(err.Error())
		}

		edgecfg, err := translate.ConfigToEdgeConfig(cfg)
		if err != nil {
			log.Fatalln(err.Error())
		}

		marshalled, err := edgeconfig.Marshal(edgecfg)
		if err != nil {
			log.Fatalln(err.Error())
		}

		// @TODO should iterate all devices that come back in edgecfg and apply
		connDeets := cfg.Routers[0].Connection
		ssh, err := connection.NewSSHConnection(connDeets.IP, connDeets.Port, connDeets.Username, connDeets.Password)
		if err != nil {
			log.Fatalln(err.Error())
		}

		live, err := ssh.FetchLiveConfig()
		if err != nil {
			log.Fatalln(err.Error())
		}

		err = os.WriteFile(fmt.Sprintf("config.boot.%d", time.Now().Unix()), live, 0644)
		if err != nil {
			log.Fatalf("error saving backup of current config: %s", err.Error())
		}

		footer := util.LastNLines(string(live), 4)
		withFooter := bytes.Join([][]byte{marshalled, []byte(footer)}, []byte("\n"))

		cfgPath := "/tmp/edgefig.cfg"
		err = ssh.WriteFile(cfgPath, withFooter)
		if err != nil {
			log.Fatalln(err.Error())
		}

		err = ssh.ApplyConfig(cfgPath)
		if err != nil {
			log.Fatalln(err.Error())
		}

		err = ssh.DeleteFile(cfgPath)
		if err != nil {
			log.Fatalln(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)
}
