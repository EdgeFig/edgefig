package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// dumpConfigCmd represents the dumpConfig command
var dumpConfigCmd = &cobra.Command{
	Use:   "dump-config",
	Short: "Dumps the generated configs to files",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("dump-config called")
	},
}

func init() {
	rootCmd.AddCommand(dumpConfigCmd)
}
