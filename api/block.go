package api

import (
	"context"
)

// Block : ユーザをブロック
func (a *API) Block(userID string) error {
	_, err := a.client.UserBlocks(context.Background(), a.CurrentUser.ID, userID)

	return checkError(err)
}

// UnBlock : ユーザのブロックを解除
func (a *API) UnBlock(userID string) error {
	_, err := a.client.DeleteUserBlocks(context.Background(), a.CurrentUser.ID, userID)

	return checkError(err)
}
