package config

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

var cfg aws.Config
var once sync.Once
var cfgError error

func Configuration() (aws.Config, error) {

	once.Do(func() {
		ctx := context.Background()
		cfg, cfgError = config.LoadDefaultConfig(ctx)
	})

	return cfg, cfgError
}
