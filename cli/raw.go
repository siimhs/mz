package cli

import (
	"os"

	"github.com/spf13/cobra"
)

// rawCmd dumps a logged-in page; used to discover and debug MZ endpoints.
var rawCmd = &cobra.Command{
	Use:    "raw <path>",
	Short:  "Print the raw response of a logged-in GET, e.g. 'mz raw \"?p=players\"'",
	Hidden: true,
	Args:   cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		body, err := c.Get(args[0])
		if err != nil {
			return err
		}
		_, err = os.Stdout.Write(body)
		return err
	},
}

func init() {
	rootCmd.AddCommand(rawCmd)
}
