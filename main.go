package main

import (
	"fmt"
	"log"

	"github.com/arrow2nd/nekome/oauth"
)

func main() {
	client := oauth.New()

	token, err := client.Auth()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(token)
}
