package api

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/g8rswimmer/go-twitter/v2"
)

func (a *API) AuthUserLookup() error {
	client, err := a.newClient()
	if err != nil {
		return err
	}

	opts := twitter.UserLookupOpts{
		Expansions: []twitter.Expansion{twitter.ExpansionPinnedTweetID},
	}

	fmt.Println("Callout to auth user lookup callout")

	userResponse, err := client.AuthUserLookup(context.Background(), opts)
	if err != nil {
		return fmt.Errorf("auth user lookup error: %v", err)
	}

	dictionaries := userResponse.Raw.UserDictionaries()

	enc, err := json.MarshalIndent(dictionaries, "", "    ")
	if err != nil {
		return err
	}

	fmt.Println(string(enc))

	return nil
}
