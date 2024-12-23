/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package category

import (
	"context"
	"fmt"

	"github.com/gurunandan-bhat/sql-to-nosql/cmd"
	"github.com/gurunandan-bhat/sql-to-nosql/internal/model"
	"github.com/gurunandan-bhat/sql-to-nosql/internal/reldb"
	"github.com/spf13/cobra"
)

// categoryCmd represents the category command
var categoryCmd = &cobra.Command{
	Use:   "category",
	Short: "Tranfer Mario Categories from mysql to dynamodb",
	RunE: func(cmd *cobra.Command, args []string) error {

		cfg, err := reldb.Configuration()
		if err != nil {
			return fmt.Errorf("error fetching configuration: %s", err)
		}

		relDBH, err := reldb.NewModel(cfg)
		if err != nil {
			return fmt.Errorf("error connecting to database: %s", err)
		}

		categories, err := relDBH.Categories()
		if err != nil {
			return fmt.Errorf("error fetching categories in cmd: %s", err)
		}

		batchSize := 200
		count, err := model.AddCategoryBatch(context.Background(), categories, batchSize)
		if err != nil {
			return fmt.Errorf("error adding bulk categories: %s", err)
		}
		fmt.Println("Inserted ", count, " categories")

		return nil
	},
}

func init() {
	cmd.RootCmd.AddCommand(categoryCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// categoryCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// categoryCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
