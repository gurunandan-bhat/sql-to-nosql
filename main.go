/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/gurunandan-bhat/sql-to-nosql/cmd"
	_ "github.com/gurunandan-bhat/sql-to-nosql/cmd/category"
	_ "github.com/gurunandan-bhat/sql-to-nosql/cmd/product"
	_ "github.com/gurunandan-bhat/sql-to-nosql/cmd/recipients"
	_ "github.com/gurunandan-bhat/sql-to-nosql/cmd/tree"
)

func main() {
	cmd.Execute()
}
