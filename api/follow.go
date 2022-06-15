package api

import (
	"context"
	"fmt"
)

// Follow ユーザをフォロー
func (a *API) Follow(userID string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	client := a.newClient(a.CurrentUser.Token)

	if _, err := client.UserFollows(context.Background(), a.CurrentUser.ID, userID); err != nil {
		return fmt.Errorf("user follow error: %v", err)
	}

	return nil
}

// UnFollow ユーザのフォローを解除
func (a *API) UnFollow(userID string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	client := a.newClient(a.CurrentUser.Token)

	if _, err := client.DeleteUserFollows(context.Background(), a.CurrentUser.ID, userID); err != nil {
		return fmt.Errorf("user unfollow error: %v", err)
	}

	return nil
}
