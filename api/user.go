package api

import (
	"context"
	"fmt"

	"github.com/g8rswimmer/go-twitter/v2"
)

// FetchUser UserNameからユーザ情報を取得
func (a *API) FetchUser(userNames []string) ([]*twitter.UserObj, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	client := a.newClient(a.CurrentUser.Token)

	opts := twitter.UserLookupOpts{
		UserFields: userFields,
	}

	result, err := client.UserNameLookup(context.Background(), userNames, opts)
	if err != nil {
		return nil, fmt.Errorf("username lookup error: %v", err)
	}

	return result.Raw.Users, nil
}
