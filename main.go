package main

import "github.com/arrow2nd/nekome/oauth"

func main() {
	client := oauth.New()

	client.Auth()
}
