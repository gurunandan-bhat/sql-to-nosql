/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package recipients

import (
	"fmt"

	"github.com/gurunandan-bhat/sql-to-nosql/cmd"
	"github.com/gurunandan-bhat/sql-to-nosql/internal/reldb"
	"github.com/spf13/cobra"
)

// recipientsCmd represents the recipients command
var recipientsCmd = &cobra.Command{
	Use:   "recipients",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		fmt.Println("recipients called")

		cfg, err := reldb.Configuration()
		if err != nil {
			return err
		}

		m, err := reldb.NewModel(cfg)
		if err != nil {
			return err
		}

		rcpts, err := m.Recipients()
		if err != nil {
			return err
		}
		fmt.Println(rcpts)

		return nil
	},
}

func init() {
	cmd.RootCmd.AddCommand(recipientsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// recipientsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// recipientsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
