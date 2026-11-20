package cli

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mz",
	Short: "This is command line interface (cli) for managerzone.com",
	Long: `Welcome to 'mz'!

This little command line tool aims to simplify interacting with managerzone.com.

* Get training reports in structured way.
* Arrange friendlies.
* Make tactic templates.
* Enrich your players with tags.

Type 'mz help' and start exploring.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// rootCmd.CompletionOptions.DisableDefaultCmd = false
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.SilenceUsage = true
}
