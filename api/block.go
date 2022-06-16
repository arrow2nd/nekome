package api

import (
	"context"
	"fmt"
)

// Block ユーザをブロック
func (a *API) Block(userID string) error {
	if _, err := a.client.UserBlocks(context.Background(), a.CurrentUser.ID, userID); err != nil {
		return fmt.Errorf("user block error: %v", err)
	}

	return nil
}

// UnBlock ユーザのブロックを解除
func (a *API) UnBlock(userID string) error {
	if _, err := a.client.DeleteUserBlocks(context.Background(), a.CurrentUser.ID, userID); err != nil {
		return fmt.Errorf("user unblock error: %v", err)
	}

	return nil
}
