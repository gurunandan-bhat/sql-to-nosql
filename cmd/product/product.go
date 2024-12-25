/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package product

import (
	"encoding/json"
	"fmt"

	"github.com/gurunandan-bhat/sql-to-nosql/cmd"
	"github.com/gurunandan-bhat/sql-to-nosql/internal/reldb"
	"github.com/spf13/cobra"
)

// productCmd represents the product command
var productCmd = &cobra.Command{
	Use:   "product",
	Short: "Transfer Mario Products from mysql to dynamodb",
	RunE: func(cmd *cobra.Command, args []string) error {

		iProdID, err := cmd.Flags().GetUint32("iProdID")
		if err != nil {
			return fmt.Errorf("error parsing argument iProdID %d: %s", iProdID, err)
		}
		fmt.Println("product called with product id", iProdID)

		cfg, err := reldb.Configuration()
		if err != nil {
			return fmt.Errorf("error fetching configuration: %s", err)
		}

		relDBH, err := reldb.NewModel(cfg)
		if err != nil {
			return fmt.Errorf("error connecting to database: %s", err)
		}

		showProducts, err := cmd.Flags().GetBool("show-products")
		if err != nil {
			return fmt.Errorf("error parsing show-products: %s", err)
		}
		if showProducts {
			products, err := relDBH.Products()
			if err != nil {
				return err
			}
			jsonProducts, err := json.MarshalIndent(&products, "", "\t")
			if err != nil {
				return fmt.Errorf("error marshalling products: %s", err)
			}
			fmt.Println("Products: ", string(jsonProducts))
			return nil
		}

		showAttribs, err := cmd.Flags().GetBool("show-attributes")
		if err != nil {
			return fmt.Errorf("error parsing show-attributes option: %s", err)
		}
		if showAttribs {
			prodAttribs, err := relDBH.ProductAttributes(iProdID)
			if err != nil {
				return fmt.Errorf("error fetching product attributes for product %d: %s", iProdID, err)
			}
			jsonBytesAttribs, err := json.MarshalIndent(&prodAttribs, "", "\t")
			if err != nil {
				return fmt.Errorf("error marshalling product attributes: %s", err)
			}
			fmt.Println(
				"Attributes: ", string(jsonBytesAttribs),
			)
		}

		showColors, err := cmd.Flags().GetBool("show-colors")
		if err != nil {
			return fmt.Errorf("error parsing show-colors option: %s", err)
		}

		if showColors {
			prodColorAttribs, err := relDBH.ProductColorAttributes(iProdID)
			if err != nil {
				return fmt.Errorf("error fetching color attributes for product %d: %s", iProdID, err)
			}
			jsonBytesColorAttribs, err := json.MarshalIndent(&prodColorAttribs, "", "\t")
			if err != nil {
				return fmt.Errorf("error marshalling product attributes: %s", err)
			}
			fmt.Println(
				"Color Attributes: ", string(jsonBytesColorAttribs),
			)
		}

		return nil
	},
}

func init() {
	cmd.RootCmd.AddCommand(productCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// productCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// productCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	productCmd.Flags().Uint32P("iProdID", "i", 0, "Display attributes of product with <id>")
	productCmd.Flags().BoolP("show-products", "p", false, "Dump products")
	productCmd.Flags().BoolP("show-attributes", "a", false, "Dump attributes")
	productCmd.Flags().BoolP("show-colors", "c", false, "Dump color attributes")
}
