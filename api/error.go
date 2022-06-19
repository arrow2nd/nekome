package api

import (
	"fmt"
	"log"

	"github.com/g8rswimmer/go-twitter/v2"
)

func checkError(err error) error {
	// レート制限
	if rateLimit, has := twitter.RateLimitFromError(err); has {
		log.Println(err)
		t := rateLimit.Reset.Time().Local().Format("15:04:05")
		return fmt.Errorf("Rate limit exceeded (Reset time: %s)", t)
	}

	return err
}
