package main

import (
	"fmt"

	"github.com/alfredo-bertorelli/blog-aggregator/internal/config"
)

func main() {
	config, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}
	config.SetUser("alfredo-bertorelli-hernandez")
	fmt.Println(config.DbURL)
	fmt.Println(config.CurrentUserName)
}
