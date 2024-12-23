package model

import (
	"context"
	"fmt"
	"log"

	"github.com/gurunandan-bhat/sql-to-nosql/internal/config"
	"github.com/gurunandan-bhat/sql-to-nosql/internal/reldb"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
)

type CategoryValue struct {
	PK string
	SK string
	reldb.CategorySummary
}

// TableBasics encapsulates the Amazon DynamoDB service actions used in the examples.
// It contains a DynamoDB service client that is used to act on the specified table.
type TableBasics struct {
	DynamoDbClient *dynamodb.Client
	TableName      string
}

func PutCategory(ctx context.Context, cat reldb.CategorySummary) error {

	cfg, err := config.Configuration()
	if err != nil {
		log.Fatalf("error fetching default configuration: %s", err)
	}

	client := dynamodb.NewFromConfig(cfg)

	catVal := CategoryValue{
		PK:              "Category",
		SK:              fmt.Sprintf("CAT#%d", cat.IPCatID),
		CategorySummary: cat,
	}

	av, err := attributevalue.MarshalMap(catVal)
	if err != nil {
		return fmt.Errorf("error converting category to attribute-value: %s", err)
	}

	_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String("Mario"),
		Item:      av,
	})
	if err != nil {
		return fmt.Errorf("error adding category: %s", err)
	}

	return nil
}

// AddMovieBatch adds a slice of movies to the DynamoDB table. The function sends
// batches of 25 movies to DynamoDB until all movies are added or it reaches the
// specified maximum.
func AddCategoryBatch(ctx context.Context, categories []reldb.CategorySummary, maxCats int) (int, error) {

	var err error
	var item map[string]types.AttributeValue

	cfg, err := config.Configuration()
	if err != nil {
		return 0, fmt.Errorf("error fetching default configuration: %s", err)
	}

	client := dynamodb.NewFromConfig(cfg)

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
				PK:              "Category",
				SK:              fmt.Sprintf("CAT#%d", category.IPCatID),
				CategorySummary: category,
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

		requests := map[string][]types.WriteRequest{"Mario": writeReqs}
		_, err = client.BatchWriteItem(ctx, &dynamodb.BatchWriteItemInput{
			RequestItems: requests,
		})
		if err != nil {
			log.Printf("Couldn't add a batch of categories to %v. Here's why: %v\n", "Mario", err)
		} else {
			written += len(writeReqs)
		}

		start = end
		end += batchSize
	}

	return written, err
}
