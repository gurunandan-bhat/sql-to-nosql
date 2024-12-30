/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/gurunandan-bhat/sql-to-nosql/cmd"
	"github.com/gurunandan-bhat/sql-to-nosql/internal/reldb"
	"github.com/spf13/cobra"
)

// treeCmd represents the tree command
var treeCmd = &cobra.Command{
	Use:   "tree",
	Short: "Generate the Category Tree",
	RunE: func(cmd *cobra.Command, args []string) error {

		fmt.Println("tree called")
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

		root := categories[0]
		for _, cat1 := range root.Children {
			fmt.Println(cat1.VName)
			for _, cat2 := range cat1.Children {
				fmt.Println("\t", cat2.VName)
				for _, cat3 := range cat2.Children {
					fmt.Println("\t\t", cat3.VName)
					for _, cat4 := range cat3.Children {
						fmt.Println("\t\t\t", cat4.VName)
					}
				}
			}
		}

		return nil
	},
}

func init() {
	cmd.RootCmd.AddCommand(treeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// treeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// treeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
