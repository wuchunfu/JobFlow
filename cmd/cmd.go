package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/wuchunfu/JobFlow/cmd/server"
	"os"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:          "job-flow",
	SilenceUsage: true,
	Short:        "Main application",
	Long:         `job-flow is a timing scheduling system implemented with golang.`,
	Example:      "job-flow job-flow",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires at least 1 arg(s), only received 0")
		}
		if cmd.Use != args[0] {
			return fmt.Errorf("invalid args specified: %s", args[0])
		}
		return nil
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Run Done.")
	},
}

func init() {
	rootCmd.AddCommand(server.StartCmd)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
