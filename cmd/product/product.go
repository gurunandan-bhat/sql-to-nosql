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

		// if iProdID == 0 {
		// 	return nil
		// }

		cfg, err := reldb.Configuration()
		if err != nil {
			return fmt.Errorf("error fetching configuration: %s", err)
		}

		relDBH, err := reldb.NewModel(cfg)
		if err != nil {
			return fmt.Errorf("error connecting to database: %s", err)
		}

		products, err := relDBH.Products()
		if err != nil {
			return err
		}

		// prodAttribs, err := relDBH.ProductAttributes(iProdID)
		// if err != nil {
		// 	return fmt.Errorf("error fetching product attributes for product %d: %s", iProdID, err)
		// }

		// prodColors, err := relDBH.ProductColors(iProdID)
		// if err != nil {
		// 	return fmt.Errorf("error fetching product colors for product %d: %s", iProdID, err)
		// }

		// prodColorAttribs, err := relDBH.ProductColorAttributes(iProdID)
		// if err != nil {
		// 	return fmt.Errorf("error fetching color attributes for product %d: %s", iProdID, err)
		// }

		jsonProducts, err := json.MarshalIndent(&products, "", "\t")
		if err != nil {
			return fmt.Errorf("error marshalling products: %s", err)
		}

		// jsonBytesAttribs, err := json.MarshalIndent(&prodAttribs, "", "\t")
		// if err != nil {
		// 	return fmt.Errorf("error marshalling product attributes: %s", err)
		// }

		// jsonBytesColors, err := json.MarshalIndent(&prodColors, "", "\t")
		// if err != nil {
		// 	return fmt.Errorf("error marshalling product colors: %s", err)
		// }

		// jsonBytesColorAttribs, err := json.MarshalIndent(&prodColorAttribs, "", "\t")
		// if err != nil {
		// 	return fmt.Errorf("error marshalling product attributes: %s", err)
		// }

		fmt.Println(
			"Products: ", string(jsonProducts),
			// "\nAttributes: ", string(jsonBytesAttribs),
			// "\nColors: ", string(jsonBytesColors),
			// "\nColor Attributes: ", string(jsonBytesColorAttribs),
		)

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
}
