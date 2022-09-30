package api

import (
	"context"
)

// Block : ユーザをブロック
func (a *API) Block(userId string) error {
	_, err := a.client.UserBlocks(context.Background(), a.CurrentUser.ID, userId)

	return checkError(err)
}

// UnBlock : ユーザのブロックを解除
func (a *API) UnBlock(userId string) error {
	_, err := a.client.DeleteUserBlocks(context.Background(), a.CurrentUser.ID, userId)

	return checkError(err)
}
