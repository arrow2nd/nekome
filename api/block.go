package api

import (
	"context"
	"fmt"
)

// Block ユーザをブロック
func (a *API) Block(userID string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	client := a.newClient(a.CurrentUser.Token)

	if _, err := client.UserBlocks(context.Background(), a.CurrentUser.ID, userID); err != nil {
		return fmt.Errorf("user block error: %v", err)
	}

	return nil
}

// UnBlock ユーザのブロックを解除
func (a *API) UnBlock(userID string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	client := a.newClient(a.CurrentUser.Token)

	if _, err := client.DeleteUserBlocks(context.Background(), a.CurrentUser.ID, userID); err != nil {
		return fmt.Errorf("user unblock error: %v", err)
	}

	return nil
}
