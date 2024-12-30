/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package category

import (
	"context"
	"encoding/json"
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

		bDryRun, err := cmd.Flags().GetBool("dry-run")
		if err != nil {
			return fmt.Errorf("error parsing argument dry-run: %s", err)
		}

		cfg, err := reldb.Configuration()
		if err != nil {
			return fmt.Errorf("error fetching configuration: %s", err)
		}

		relDBH, err := reldb.NewModel(cfg)
		if err != nil {
			return fmt.Errorf("error connecting to database: %s", err)
		}

		categories, err := relDBH.CategoryTree()
		if err != nil {
			return fmt.Errorf("error fetching categories in cmd: %s", err)
		}

		if !bDryRun {
			batchSize := 200
			count, err := model.AddCategoryBatch(context.Background(), categories, batchSize)
			if err != nil {
				return fmt.Errorf("error adding bulk categories: %s", err)
			}
			fmt.Println("Inserted ", count, " categories")
			return nil
		}

		jsonBytes, err := json.MarshalIndent(&categories, "", "\t")
		if err != nil {
			return fmt.Errorf("error marshaling categories: %s", err)
		}
		fmt.Println("Categories: ", string(jsonBytes))

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
	categoryCmd.Flags().BoolP("dry-run", "d", false, "Dump categories, dont insert")
}
