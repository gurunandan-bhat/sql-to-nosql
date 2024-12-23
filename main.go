/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/gurunandan-bhat/sql-to-nosql/cmd"
	_ "github.com/gurunandan-bhat/sql-to-nosql/cmd/category"
	_ "github.com/gurunandan-bhat/sql-to-nosql/cmd/product"
)

func main() {
	cmd.Execute()
}
