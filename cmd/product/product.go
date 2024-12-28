/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package product

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gurunandan-bhat/sql-to-nosql/cmd"
	"github.com/gurunandan-bhat/sql-to-nosql/internal/model"
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
			return fmt.Errorf("error parsing show-attributes: %s", err)
		}
		if showAttribs {
			attribs, err := relDBH.ProductAttributes(iProdID, false)
			if err != nil {
				return fmt.Errorf("error fetching product attributes for %d: %s", iProdID, err)
			}
			jsonAttribs, err := json.MarshalIndent(&attribs, "", "\t")
			if err != nil {
				return fmt.Errorf("error marshalling attributes: %s", err)
			}
			fmt.Println("Attributes: ", string(jsonAttribs))
			return nil
		}

		showSKUs, err := cmd.Flags().GetBool("show-skus")
		if err != nil {
			return fmt.Errorf("error parsing show-skus option: %s", err)
		}
		if showSKUs {
			productSKUs, err := relDBH.ProductSKUs(iProdID)
			if err != nil {
				return fmt.Errorf("error fetching product skus for product %d: %s", iProdID, err)
			}
			jsonBytesAttribs, err := json.MarshalIndent(&productSKUs, "", "\t")
			if err != nil {
				return fmt.Errorf("error marshalling product attributes: %s", err)
			}
			fmt.Println(
				"SKUs: ", string(jsonBytesAttribs),
			)
		}

		products, err := relDBH.Products()
		if err != nil {
			return fmt.Errorf("error fetching all products: %s", err)
		}

		for i := range products {

			iProdID := products[i].IProdID
			attribs, err := relDBH.ProductAttributes(iProdID, false)
			if err != nil {
				return fmt.Errorf("error fetching product attributes for %d: %s", iProdID, err)
			}
			products[i].Attributes = attribs

			skus, err := relDBH.ProductSKUs(iProdID)
			if err != nil {
				return fmt.Errorf("error fetching product skus for product %d: %s", iProdID, err)
			}
			products[i].SKUs = skus
		}

		batchSize := 3000
		count, err := model.AddProductBatch(context.TODO(), products, batchSize)
		if err != nil {
			return fmt.Errorf("error adding products to dynamodb: %s", err)
		}

		fmt.Printf("Inserted %d products\n", count)

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
	productCmd.Flags().BoolP("show-skus", "s", false, "Dump product SKUs")
	productCmd.Flags().BoolP("show-attributes", "a", false, "Dump product attributes")
}
