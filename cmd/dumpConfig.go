package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/cmmarslender/edgefig/pkg/config"
	"github.com/cmmarslender/edgefig/pkg/edgeconfig"
	"github.com/cmmarslender/edgefig/pkg/translate"
)

// dumpConfigCmd represents the dumpConfig command
var dumpConfigCmd = &cobra.Command{
	Use:   "dump-config",
	Short: "Dumps the generated configs to files",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig("config.yml")
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

		err = os.WriteFile("config-out", marshalled, 0644)
		if err != nil {
			log.Fatalln(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(dumpConfigCmd)
}
