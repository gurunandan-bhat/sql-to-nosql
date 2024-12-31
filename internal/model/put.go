package model

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gurunandan-bhat/sql-to-nosql/internal/config"
	"github.com/gurunandan-bhat/sql-to-nosql/internal/reldb"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type CategoryValue struct {
	PK                string
	SK                string
	IPCatID           uint32
	VCategoryName     string
	VCategoryURLName  string
	IParentID         uint32
	VShortDescription *string
	MImages           reldb.Images
	CTypeStatus       string
	LAttributes       []reldb.CategoryAttribute
	LChildren         []*reldb.CategorySummary
}

type ProductValue struct {
	PK                string
	SK                string
	IProdID           uint32
	IPCatID           uint32
	CCode             *string
	VName             string
	VURLName          string
	VCategoryName     string
	VCategoryURLName  string
	VShortDescription *string
	VDescription      *string
	MPrices           reldb.ProdPrice
	MImages           reldb.Images
	CTypeStatus       string
	VYTID             *string
	LAttributes       []reldb.ProductAttribute
	LSKUs             []reldb.SKU
}

// TableBasics encapsulates the Amazon DynamoDB service actions used in the examples.
// It contains a DynamoDB service client that is used to act on the specified table.
type TableBasics struct {
	DynamoDbClient *dynamodb.Client
	TableName      string
}

var tableName = "MarioGallery"
var sleepBetweenBatches time.Duration = 30 * time.Second

func PutCategory(ctx context.Context, cat reldb.CategorySummary) error {

	cfg, err := config.Configuration()
	if err != nil {
		log.Fatalf("error fetching default configuration: %s", err)
	}

	// client := dynamodb.NewFromConfig(cfg)

	client := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = aws.String("https://localhost:4000")
	})

	catVal := CategoryValue{
		PK:                "Category",
		SK:                fmt.Sprintf("CAT#%d", cat.IPCatID),
		IPCatID:           cat.IPCatID,
		VCategoryName:     cat.VName,
		VCategoryURLName:  cat.VURLName,
		IParentID:         cat.IParentID,
		VShortDescription: cat.VShortDesc,
		MImages:           cat.Images,
		CTypeStatus:       fmt.Sprintf("C%s", cat.CStatus),
		LAttributes:       cat.Attributes,
		LChildren:         cat.Children,
	}

	av, err := attributevalue.MarshalMap(catVal)
	if err != nil {
		return fmt.Errorf("error converting category to attribute-value: %s", err)
	}

	_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      av,
	})
	if err != nil {
		return fmt.Errorf("error adding category: %s", err)
	}

	return nil
}

// AddCategoryBatch adds a slice of categories to the DynamoDB table. The function sends
// batches of 25 categories to DynamoDB until all categories are added or it reaches the
// specified maximum.
func AddCategoryBatch(ctx context.Context, categories []reldb.CategorySummary, maxCats int) (int, error) {

	var err error
	var item map[string]types.AttributeValue

	cfg, err := config.Configuration()
	if err != nil {
		return 0, fmt.Errorf("error fetching default configuration: %s", err)
	}

	client := dynamodb.NewFromConfig(cfg)

	// client := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
	// 	o.BaseEndpoint = aws.String("http://localhost:4000")
	// })

	written := 0
	batchSize := 25 // DynamoDB allows a maximum batch size of 25 items.
	start := 0
	end := start + batchSize
	for start < maxCats && start < len(categories) {

		var writeReqs []types.WriteRequest
		if end > len(categories) {
			end = len(categories)
		}

		for _, category := range categories[start:end] {
			catVal := CategoryValue{
				PK:                "Category",
				SK:                fmt.Sprintf("CAT#%d", category.IPCatID),
				IPCatID:           category.IPCatID,
				VCategoryName:     category.VName,
				VCategoryURLName:  category.VURLName,
				IParentID:         category.IParentID,
				VShortDescription: category.VShortDesc,
				MImages:           category.Images,
				CTypeStatus:       fmt.Sprintf("C%s", category.CStatus),
				LAttributes:       category.Attributes,
				LChildren:         category.Children,
			}

			item, err = attributevalue.MarshalMap(catVal)
			if err != nil {
				log.Printf("Couldn't marshal category %v for batch writing. Here's why: %v\n", category.VName, err)
			} else {
				writeReqs = append(
					writeReqs,
					types.WriteRequest{PutRequest: &types.PutRequest{Item: item}},
				)
			}
		}

		requests := map[string][]types.WriteRequest{tableName: writeReqs}
		_, err = client.BatchWriteItem(ctx, &dynamodb.BatchWriteItemInput{
			RequestItems: requests,
		})
		if err != nil {
			log.Printf("Couldn't add a batch of categories to %v. Here's why: %v\n", tableName, err)
		} else {
			written += len(writeReqs)
		}

		start = end
		end += batchSize

		fmt.Printf("Sleeping for %s after adding %d entries\n", sleepBetweenBatches, start)
		time.Sleep(sleepBetweenBatches)

	}

	return written, err
}

// AddCategoryBatch adds a slice of categories to the DynamoDB table. The function sends
// batches of 25 categories to DynamoDB until all categories are added or it reaches the
// specified maximum.
func AddProductBatch(ctx context.Context, products []reldb.Product, maxCats int) (int, error) {

	cfg, err := config.Configuration()
	if err != nil {
		return 0, fmt.Errorf("error fetching default configuration: %s", err)
	}

	client := dynamodb.NewFromConfig(cfg)

	// client := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
	// 	o.BaseEndpoint = aws.String("http://localhost:4000")
	// })

	written := 0
	batchSize := 25 // DynamoDB allows a maximum batch size of 25 items.
	start := 0
	end := start + batchSize
	for start < maxCats && start < len(products) {

		var writeReqs []types.WriteRequest
		if end > len(products) {
			end = len(products)
		}

		for _, product := range products[start:end] {
			prodVal := ProductValue{
				PK:                fmt.Sprintf("%sPROD%d", *product.CStatus, product.IProdID),
				SK:                fmt.Sprintf("CAT#%d", product.IPCatID),
				IPCatID:           product.IPCatID,
				VName:             product.VName,
				VURLName:          product.VURLName,
				VCategoryName:     product.VCategoryName,
				VCategoryURLName:  product.VCategoryURLName,
				CCode:             product.CCode,
				VShortDescription: product.VShortDesc,
				VDescription:      product.VDescription,
				MImages:           product.Images,
				MPrices:           product.ProdPrice,
				CTypeStatus:       fmt.Sprintf("P%s", *product.CStatus),
				LAttributes:       product.Attributes,
				LSKUs:             product.SKUs,
				VYTID:             product.VYTID,
			}

			item, err := attributevalue.MarshalMap(prodVal)
			if err != nil {
				log.Printf("Couldn't marshal category %v for batch writing. Here's why: %v\n", product.VName, err)
			} else {
				writeReqs = append(
					writeReqs,
					types.WriteRequest{PutRequest: &types.PutRequest{Item: item}},
				)
			}
		}

		requests := map[string][]types.WriteRequest{tableName: writeReqs}
		_, err := client.BatchWriteItem(ctx, &dynamodb.BatchWriteItemInput{
			RequestItems: requests,
		})
		if err != nil {
			log.Printf("Couldn't add a batch of products to %v. Here's why: %v\n", tableName, err)
		} else {
			written += len(writeReqs)
		}

		start = end
		end += batchSize

		fmt.Printf("Sleeping for %s after adding %d entries\n", sleepBetweenBatches, start)
		time.Sleep(sleepBetweenBatches)
	}

	return written, err
}
